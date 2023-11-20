package plugins

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	v1 "github.com/neermitt/opsos/api/v1"
	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/stack"
	"github.com/neermitt/opsos/pkg/utils"
	"github.com/spf13/afero"
	"io"
	"log"
	"os"
)

type KubeConfigProvider interface {
	ExportKubeConfig(ctx context.Context, stk *stack.Stack, providerStackSettings map[string]any, kubeConfigPath string) error
}

type KubeConfigProviderFactory func(conf *v1.ConfigSpec) (KubeConfigProvider, error)

var kubeConfigProviderFactoryMap map[string]KubeConfigProviderFactory

func init() {
	kubeConfigProviderFactoryMap = make(map[string]KubeConfigProviderFactory)
}

func RegisterKubeConfigProvider(name string, factory KubeConfigProviderFactory) {
	kubeConfigProviderFactoryMap[name] = factory
}

func GetKubeConfigProvider(ctx context.Context, name string) (KubeConfigProvider, error) {
	pf, ok := kubeConfigProviderFactoryMap[name]
	if !ok {
		return nil, fmt.Errorf("kubeconfig provider `%s` not found", name)
	}
	return pf(config.GetConfig(ctx))
}

func GetKubeConfig(ctx context.Context, k8sSettings map[string]any, stk *stack.Stack, componentName string, kubeConfigPath string) error {
	log.Printf("[TRACE] Begin loading kubeconfig for stack %[1]s", stk.Id)
	defer func() {
		log.Printf("[TRACE] End loading kubeconfig for stack %[1]s", stk.Id)
	}()
	provider, found := k8sSettings["provider"].(string)
	if !found {
		msg := fmt.Sprintf("k8s settings for provider %[3]s not found for component %[2]s in stack %[1]s", stk.Id, componentName, provider)
		log.Printf("[ERROR] %s", msg)
		return errors.New(msg)
	}
	log.Printf("[INFO] loading kubeconfig for stack %[1]s using provider %[2]s", stk.Id, provider)
	kubeConfigProvider, err := GetKubeConfigProvider(ctx, provider)
	if err != nil {
		return err
	}
	providerSettings := k8sSettings[provider].(map[string]any)

	log.Printf("[INFO] Writing the kubeconfig to file: %[1]s", kubeConfigPath)

	return kubeConfigProvider.ExportKubeConfig(ctx, stk, providerSettings, kubeConfigPath)
}

type kubeFileProviderSettings struct {
	Path string `yaml:"path" json:"path" mapstructure:"path" validate:"required"`
}

type kubeConfigFileProvider struct {
}

func (k kubeConfigFileProvider) ExportKubeConfig(_ context.Context, _ *stack.Stack, providerStackSettings map[string]any, kubeConfigPath string) error {
	var settings kubeFileProviderSettings
	err := utils.FromMap(providerStackSettings, &settings)

	validate := validator.New()
	if err = validate.Struct(&settings); err != nil {
		return err
	}

	data, err := afero.ReadFile(afero.NewOsFs(), settings.Path)
	if err != nil {
		return err
	}

	return exportKindKubeConfigRaw(data, kubeConfigPath)
}

func exportKindKubeConfigRaw(kubeConfig []byte, kubeConfigPath string) (err error) {
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

	_, err = io.Copy(f, bytes.NewReader(kubeConfig))
	return err
}

func init() {
	RegisterKubeConfigProvider("file", func(conf *v1.ConfigSpec) (KubeConfigProvider, error) {
		return kubeConfigFileProvider{}, nil
	})
}
