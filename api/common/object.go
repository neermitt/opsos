package common

type ObjectMetadata struct {
	Name        string `yaml:"name" json:"name" mapstructure:"name"`
	Description string `yaml:"description" json:"description" mapstructure:"description"`
}

type Object struct {
	ApiVersion string         `yaml:"apiVersion" json:"apiVersion" mapstructure:"apiVersion"`
	Kind       string         `yaml:"kind" json:"kind" mapstructure:"kind"`
	Metadata   ObjectMetadata `yaml:"metadata" json:"metadata"`
}
