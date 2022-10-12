package schema

import "github.com/neermitt/opsos/pkg/components"

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
	Types map[string]map[string]components.ConfigWithMetadata `yaml:",inline" json:",inline" mapstructure:",remain"`
}
