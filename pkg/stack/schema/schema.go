package schema

type StackConfig struct {
	Vars                  map[string]any                   `yaml:"vars,omitempty" json:"vars,omitempty" mapstructure:"vars"`
	Envs                  map[string]string                `yaml:"envs,omitempty" json:"envs,omitempty" mapstructure:"envs"`
	KubeConfigProvider    string                           `yaml:"kube_config_provider,omitempty" json:"kube_config_provider,omitempty" mapstructure:"kube_config_provider"`
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
}

type ComponentsConfig struct {
	Types map[string]map[string]ComponentConfig `yaml:",inline" json:",inline" mapstructure:",remain"`
}

type ComponentConfig struct {
	Vars                   map[string]any    `yaml:"vars,omitempty" json:"vars,omitempty" mapstructure:"vars"`
	Envs                   map[string]string `yaml:"envs,omitempty" json:"envs,omitempty" mapstructure:"envs"`
	Component              string            `yaml:"component,omitempty" json:"component,omitempty" mapstructure:"component"`
	BackendType            string            `yaml:"backend_type,omitempty" json:"backend_type,omitempty" mapstructure:"backend_type"`
	Backend                map[string]any    `yaml:"backend,omitempty" json:"backend,omitempty" mapstructure:"backend"`
	RemoteStateBackendType string            `yaml:"remote_state_backend_type,omitempty" json:"remote_state_backend_type,omitempty" mapstructure:"remote_state_backend_type"`
	RemoteStateBackend     map[string]any    `yaml:"remote_state_backend,omitempty" json:"remote_state_backend,omitempty" mapstructure:"remote_state_backend"`
	Metadata               map[string]any    `yaml:"metadata,omitempty" json:"metadata,omitempty" mapstructure:"metadata"`
	Settings               map[string]any    `yaml:"settings,omitempty" json:"settings,omitempty" mapstructure:"settings"`
}
