package terraform

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-playground/validator"
	v1 "github.com/neermitt/opsos/api/v1"
	"github.com/neermitt/opsos/pkg/components"
	"github.com/neermitt/opsos/pkg/plugins"
	"github.com/neermitt/opsos/pkg/stack"
	"github.com/neermitt/opsos/pkg/utils"
	"github.com/variantdev/vals"
)

type kubeProviderSettings struct {
	Component string `yaml:"component" json:"component" mapstructure:"component" validate:"required"`
	Ref       string `yaml:"ref" json:"ref" mapstructure:"ref" validate:"required"`
}

func NewKubeConfigProvider(conf *v1.ConfigSpec) (plugins.KubeConfigProvider, error) {
	return &kubeConfigProvider{config: conf}, nil
}

func init() {
	plugins.RegisterKubeConfigProvider(ComponentType, NewKubeConfigProvider)
}

type kubeConfigProvider struct {
	config            *v1.ConfigSpec
	baseComponentPath string
}

func (p *kubeConfigProvider) ExportKubeConfig(ctx context.Context, stk *stack.Stack, providerStackSettings map[string]any, kubeConfigPath string) error {
	var settings kubeProviderSettings
	err := utils.FromMap(providerStackSettings, &settings)

	validate := validator.New()
	if err = validate.Struct(&settings); err != nil {
		return err
	}

	sourceComponentConfig, err := loadComponent(ctx, stk, settings.Component)
	if err != nil {
		return err
	}

	workingDir := components.GetWorkingDirectory(p.config, ComponentType, sourceComponentConfig.Component)

	workspaceName, err := ConstructWorkspaceName(stk, settings.Component, sourceComponentConfig)

	kubeConfig, err := fetchKubeConfigFromTfState(workingDir, workspaceName, settings)
	if err != nil {
		return err
	}
	return exportKindKubeConfigRaw(kubeConfig, kubeConfigPath)
}

func fetchKubeConfigFromTfState(workingDir string, workspaceName string, settings kubeProviderSettings) (string, error) {
	runtime, err := vals.New(vals.Options{
		CacheSize:     256,
		ExcludeSecret: true,
	})
	if err != nil {
		return "", err
	}

	valsRendered, err := runtime.Eval(map[string]interface{}{
		"kubeconfig": fmt.Sprintf("ref+tfstate://%s/terraform.tfstate.d/%s/terraform.tfstate/%s", workingDir, workspaceName, settings.Ref),
	})
	if err != nil {
		return "", err
	}
	kubeConfig := valsRendered["kubeconfig"].(string)
	return kubeConfig, nil
}

func loadComponent(ctx context.Context, stk *stack.Stack, componentName string) (stack.ConfigWithMetadata, error) {
	component := stack.Component{Type: ComponentType, Name: componentName}

	ctx = stack.SetStackName(ctx, stk.Id)
	ctx = stack.SetComponent(ctx, component)
	componentStack, err := stack.LoadStack(ctx, stack.LoadStackOptions{Stack: stk.Id, Component: &component})
	if err != nil {
		return stack.ConfigWithMetadata{}, err
	}
	sourceComponentConfig, found := componentStack.Components[ComponentType][componentName]
	if !found {
		return stack.ConfigWithMetadata{}, fmt.Errorf("k8s source component %s not found in stack %s", componentName, stk.Id)
	}
	return sourceComponentConfig, nil
}

func exportKindKubeConfigRaw(kubeConfig string, kubeConfigPath string) (err error) {
	f, err := os.OpenFile(kubeConfigPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := f.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	_, err = io.Copy(f, strings.NewReader(kubeConfig))
	return err
}
