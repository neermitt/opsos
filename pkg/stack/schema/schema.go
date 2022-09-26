package schema

type StackConfig struct {
	Vars                  map[string]any                   `yaml:"vars,omitempty" json:"vars,omitempty" mapstructure:"vars"`
	Components            ComponentsConfig                 `yaml:"components,omitempty" json:"components,omitempty" mapstructure:"components"`
	ComponentTypeSettings map[string]ComponentTypeSettings `yaml:",inline" json:",inline" mapstructure:",remain"`
}

type ComponentTypeSettings struct {
	Vars map[string]any `yaml:"vars,omitempty" json:"vars,omitempty" mapstructure:"vars"`
}

type ComponentsConfig struct {
	Types map[string]map[string]ComponentConfig `yaml:",inline" json:",inline" mapstructure:",remain"`
}

type ComponentConfig struct {
	Vars      map[string]any `yaml:"vars,omitempty" json:"vars,omitempty" mapstructure:"vars"`
	Component string         `yaml:"component,omitempty" json:"component,omitempty" mapstructure:"component"`
}
