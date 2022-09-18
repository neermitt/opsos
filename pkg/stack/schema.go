package stack

type Vars struct {
	Vars map[string]any `yaml:"vars,omitempty" json:"vars,omitempty"`
}

type StackConfig struct {
	Vars
	Helmfile   Vars             `yaml:"helmfile,omitempty" json:"helmfile,omitempty"`
	Components ComponentsConfig `yaml:"components,omitempty" json:"components,omitempty"`
}

type ComponentsConfig struct {
	Helmfile map[string]HelmfileConfig `yaml:"helmfile,omitempty" json:"helmfile,omitempty"`
}

type HelmfileConfig struct {
	Vars
	Component string `yaml:"component,omitempty" json:"component,omitempty"`
}
