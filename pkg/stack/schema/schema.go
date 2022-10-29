package schema

import "gopkg.in/yaml.v3"

type StackConfig struct {
	Vars                  map[string]any                   `yaml:"vars,omitempty" json:"vars,omitempty" mapstructure:"vars"`
	Envs                  map[string]string                `yaml:"env,omitempty" json:"env,omitempty" mapstructure:"env"`
	Settings              map[string]any                   `yaml:"settings,omitempty" json:"settings,omitempty" mapstructure:"settings"`
	Components            ComponentsConfig                 `yaml:"components,omitempty" json:"components,omitempty" mapstructure:"components"`
	ComponentTypeSettings map[string]ComponentTypeSettings `yaml:",inline" json:",inline" mapstructure:",remain"`
}

type ComponentTypeSettings struct {
	Vars                   map[string]any    `yaml:"vars,omitempty" json:"vars,omitempty" mapstructure:"vars"`
	Envs                   map[string]string `yaml:"envs,omitempty" json:"envs,omitempty" mapstructure:"envs"`
	BackendType            string            `yaml:"backend_type,omitempty" json:"backend_type,omitempty" mapstructure:"backend_type"`
	Backend                map[string]any    `yaml:"backend,omitempty" json:"backend,omitempty" mapstructure:"backend"`
	RemoteStateBackendType string            `yaml:"remote_state_backend_type,omitempty" json:"remote_state_backend_type,omitempty" mapstructure:"remote_state_backend_type"`
	RemoteStateBackend     map[string]any    `yaml:"remote_state_backend,omitempty" json:"remote_state_backend,omitempty" mapstructure:"remote_state_backend"`
	Settings               map[string]any    `yaml:"settings,omitempty" json:"settings,omitempty" mapstructure:"settings"`
}

type ComponentsConfig struct {
	Types map[string]map[string]ConfigWithMetadata `yaml:",inline" json:",inline" mapstructure:",remain"`
}

type ConfigWithMetadata struct {
	Config   `yaml:",inline" json:",inline" mapstructure:",squash"`
	Metadata *Metadata `yaml:"metadata,omitempty" json:"metadata,omitempty" mapstructure:"metadata,omitempty"`
}

func NewConfigFromMap(config map[string]any) (Config, error) {
	yamlCurrent, err := yaml.Marshal(config)
	if err != nil {
		return Config{}, err
	}

	var dataCurrent Config
	if err = yaml.Unmarshal(yamlCurrent, &dataCurrent); err != nil {
		return Config{}, err
	}

	return dataCurrent, nil
}

type Config struct {
	Command                   *string           `yaml:"command,omitempty" json:"command,omitempty" mapstructure:"command,omitempty"`
	Component                 *string           `yaml:"component,omitempty" json:"component,omitempty" mapstructure:"component,omitempty"`
	Vars                      map[string]any    `yaml:"vars,omitempty" json:"vars,omitempty"  mapstructure:"vars,omitempty"`
	Envs                      map[string]string `yaml:"env,omitempty" json:"env,omitempty"  mapstructure:"env,omitempty"`
	BackendType               *string           `yaml:"backend_type,omitempty" json:"backend_type,omitempty"  mapstructure:"backend_type,omitempty"`
	BackendConfigs            map[string]any    `yaml:"backend,omitempty" json:"backend,omitempty"  mapstructure:"backend,omitempty"`
	RemoteStateBackendType    *string           `yaml:"remote_state_backend_type,omitempty" json:"remote_state_backend_type,omitempty"  mapstructure:"remote_state_backend_type,omitempty"`
	RemoteStateBackendConfigs map[string]any    `yaml:"remote_state_backend,omitempty" json:"remote_state_backend,omitempty" mapstructure:"remote_state_backend,omitempty"`
	Settings                  map[string]any    `yaml:"settings,omitempty" json:"settings,omitempty" mapstructure:"settings,omitempty"`
}

func (c Config) ToMap() (map[string]any, error) {
	yamlCurrent, err := yaml.Marshal(c)
	if err != nil {
		return nil, err
	}

	var dataCurrent map[string]any
	if err = yaml.Unmarshal(yamlCurrent, &dataCurrent); err != nil {
		return nil, err
	}

	return dataCurrent, nil
}

type Metadata struct {
	Type                      *string  `yaml:"type,omitempty" json:"type,omitempty" mapstructure:"type,omitempty"`
	Component                 *string  `yaml:"component,omitempty" json:"component,omitempty" mapstructure:"component,omitempty"`
	Inherits                  []string `yaml:"inherits,omitempty" json:"inherits,omitempty" mapstructure:"inherits,omitempty"`
	TerraformWorkspace        *string  `yaml:"terraform_workspace,omitempty" json:"terraform_workspace,omitempty" mapstructure:"terraform_workspace,omitempty"`
	TerraformWorkspacePattern *string  `yaml:"terraform_workspace_pattern,omitempty" json:"terraform_workspace_pattern,omitempty" mapstructure:"terraform_workspace_pattern,omitempty"`
}
