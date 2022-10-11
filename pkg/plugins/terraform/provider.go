package terraform

import (
	"context"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"path"
	"path/filepath"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/merge"
	"github.com/neermitt/opsos/pkg/plugins"
	"github.com/neermitt/opsos/pkg/stack/schema"
	"github.com/spf13/afero"
)

const (
	ComponentType = "terraform"

	terraformStackConfigKey = "terraform-stack-config"
	stackConfigKey          = "stack-config"

	componentField              = "component"
	varsField                   = "vars"
	envsField                   = "envs"
	backendTypeField            = "backend_type"
	backendField                = "backend"
	remoteStateBackendTypeField = "remote_state_backend_type"
	remoteStateBackendField     = "remote_state_backend"
	settingsField               = "settings"
	metadataField               = "metadata"
)

var allComponentFields = []string{varsField, envsField, backendTypeField, backendField, remoteStateBackendTypeField, remoteStateBackendField, settingsField, metadataField}

type ComponentInfo struct {
	Component             string `yaml:"-" json:"-"`
	ComponentFolderPrefix string `yaml:"-" json:"-"`

	FinalComponent string            `yaml:"FinalComponent,omitempty" json:"FinalComponent,omitempty"`
	Vars           map[string]any    `yaml:"vars" json:"vars"`
	Envs           map[string]string `yaml:"envs" json:"envs"`
	BackendType    string            `yaml:"backend_type" json:"backend_type"`
	Backend        map[string]any    `yaml:"backend" json:"backend"`
}

func NewProvider(conf *config.Configuration) *Provider {
	terraformBasePath := path.Join(conf.BasePath, conf.Components.Terraform.BasePath)
	terraformBaseAbsPath, err := filepath.Abs(terraformBasePath)
	if err != nil {
		return nil
	}

	stackFS := afero.NewBasePathFs(afero.NewOsFs(), terraformBaseAbsPath)

	return &Provider{fs: stackFS, basePath: terraformBaseAbsPath}
}

type Provider struct {
	fs       afero.Fs
	basePath string
}

func (p *Provider) InitStackContext(ctx context.Context, config schema.StackConfig) context.Context {
	terraformSettings := config.ComponentTypeSettings[ComponentType]
	globalSettings := map[string]any{
		varsField: config.Vars,
		envsField: config.Envs,
	}
	stackTerraformettings := map[string]any{
		varsField:                   terraformSettings.Vars,
		envsField:                   terraformSettings.Envs,
		backendTypeField:            terraformSettings.BackendType,
		backendField:                terraformSettings.Backend,
		remoteStateBackendTypeField: terraformSettings.RemoteStateBackendType,
		remoteStateBackendField:     terraformSettings.RemoteStateBackend,
	}
	terraformSettingsOvr, _ := merge.Merge([]map[string]any{globalSettings, stackTerraformettings})
	return context.WithValue(
		context.WithValue(ctx, stackConfigKey, &config),
		terraformStackConfigKey, terraformSettingsOvr)
}

func (p *Provider) ProcessComponentConfig(ctx context.Context, componentName string, _ schema.ComponentConfig) (any, error) {
	componentConfig, err := processComponent(ctx, componentName, p.getComponentWorkingDir)

	var backendType string
	var backendMap map[string]any
	var found bool
	// select backend based on backendType
	if backendMap, found = componentConfig[backendField].(map[string]any); found {
		backendType = componentConfig[backendTypeField].(string)
		if backend, found := backendMap[backendType]; found {
			componentConfig[backendField] = backend
		} else {
			return nil, fmt.Errorf("no settings found for backend type: %s", backendType)
		}

	} else {
		return nil, errors.New("invalid `backend` field, should be a map")
	}

	// select remote state backend based on remote state backend type or backend type
	remoteStateBackendMap := backendMap

	if val, found := componentConfig[remoteStateBackendField].(map[string]any); found {
		var err error
		if remoteStateBackendMap, err = merge.Merge([]map[string]any{remoteStateBackendMap, val}); err != nil {
			return nil, err
		}
	}

	var remoteStateBackendType string
	if val, found := componentConfig[remoteStateBackendTypeField].(string); found && val != "" {
		remoteStateBackendType = val
	} else {
		remoteStateBackendType = backendType
		componentConfig[remoteStateBackendTypeField] = backendType
	}
	if remoteStateBackend, found := remoteStateBackendMap[remoteStateBackendType]; found {
		componentConfig[remoteStateBackendField] = remoteStateBackend
	} else {
		return nil, fmt.Errorf("no settings found for remote state backend type: %s", remoteStateBackendType)
	}

	// Override component if metadata has it
	if metadata, found := componentConfig[metadataField].(map[string]any); found {
		if component, found := metadata["component"]; found {
			componentConfig[componentField] = component
		}
	}

	return componentConfig, err
}

func (p *Provider) getComponentWorkingDir(componentName string) (string, error) {
	compoDirInfo, err := p.fs.Stat(componentName)
	if err != nil {
		return "", err
	}
	if !compoDirInfo.IsDir() {
		return "", fmt.Errorf("component directory %s not found in base path %s", componentName, p.basePath)
	}
	return path.Join(p.basePath, componentName), nil
}

type componentPathProvider func(componentName string) (string, error)

func processComponent(ctx context.Context, componentName string, componentPathFunc componentPathProvider) (map[string]any, error) {
	stackConfig := ctx.Value(stackConfigKey).(*schema.StackConfig)
	componentConfig, found := stackConfig.Components.Types[ComponentType][componentName]
	if !found {
		return nil, fmt.Errorf("component %s not found", componentName)
	}

	componentInfo := make(map[string]any, len(allComponentFields))

	if componentConfig.Component != "" {
		var err error
		baseComponentConfig, err := processComponent(ctx, componentConfig.Component, componentPathFunc)
		if err != nil {
			return nil, err
		}

		if err = mergeComponentConfig(&componentInfo, baseComponentConfig, componentConfig); err != nil {
			return nil, err
		}

		return componentInfo, nil
	}

	terraformStackConfig := ctx.Value(terraformStackConfigKey).(map[string]any)

	if err := mergeComponentConfig(&componentInfo, terraformStackConfig, componentConfig); err != nil {
		return nil, err
	}

	return componentInfo, nil
}

func mergeComponentConfig(dst *map[string]any, baseConfig map[string]any, componentConfig schema.ComponentConfig) error {
	componentConf := make(map[string]any, len(allComponentFields))

	if componentConfig.Vars != nil {
		componentConf[varsField] = componentConfig.Vars
	}
	if componentConfig.Envs != nil {
		componentConf[envsField] = componentConfig.Envs
	}

	if componentConfig.BackendType != "" {
		componentConf[backendTypeField] = componentConfig.BackendType
	}
	if componentConfig.Backend != nil {
		componentConf[backendField] = componentConfig.Backend
	}
	if componentConfig.RemoteStateBackendType != "" {
		componentConf[remoteStateBackendTypeField] = componentConfig.RemoteStateBackendType
	}
	if componentConfig.RemoteStateBackend != nil {
		componentConf[remoteStateBackendField] = componentConfig.RemoteStateBackend
	}
	if componentConfig.Metadata != nil {
		componentConf[metadataField] = componentConfig.Metadata
	}
	if componentConfig.Settings != nil {
		componentConf[settingsField] = componentConfig.Settings
	}

	componentOvrConfig, err := merge.Merge([]map[string]any{baseConfig, componentConf})
	if err != nil {
		return err
	}

	yamlCurrent, err := yaml.Marshal(componentOvrConfig)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(yamlCurrent, dst); err != nil {
		return err
	}

	return nil
}

func init() {
	plugins.RegisterProvider(ComponentType, func(conf *config.Configuration) plugins.Provider {
		return NewProvider(conf)
	})
}
