package v1

import (
	"github.com/neermitt/opsos/api/common"
)

type StacksSpec struct {
	BasePath      *string  `yaml:"base_path,omitempty" json:"base_path,omitempty" mapstructure:"base_path" validate:"required"`
	IncludedPaths []string `yaml:"included_paths,omitempty" json:"included_paths,omitempty" mapstructure:"included_paths" validate:"required"`
	ExcludedPaths []string `yaml:"excluded_paths,omitempty" json:"excluded_paths,omitempty" mapstructure:"excluded_paths"`
	NamePattern   *string  `yaml:"name_pattern,omitempty" json:"name_pattern,omitempty" mapstructure:"name_pattern" validate:"required"`
}

type WorkflowsSpec struct {
}

type LogSpec struct {
	Level *string `yaml:"level" json:"level" mapstructure:"level"`
	JSON  bool    `yaml:"json" json:"json" mapstructure:"json"`
	File  *string `yaml:"file" json:"file" mapstructure:"file"`
}

type ProviderSettings map[string]any

type ConfigSpec struct {
	BasePath  *string                     `yaml:"base_path,omitempty" json:"base_path,omitempty" mapstructure:"base_path" validate:"required"`
	Stacks    *StacksSpec                 `yaml:"stacks,omitempty" json:"stacks,omitempty" mapstructure:"stacks" validate:"required"`
	Workflows WorkflowsSpec               `yaml:"workflows,omitempty" json:"workflows,omitempty"`
	Logs      LogSpec                     `yaml:"logs" json:"logs" mapstructure:"logs" validate:"required"`
	Providers map[string]ProviderSettings `yaml:",inline" json:",inline" mapstructure:",remain"`
}

type Config struct {
	common.Object `yaml:",inline" json:",inline"`
	Spec          ConfigSpec `yaml:"spec" json:"spec"`
}
