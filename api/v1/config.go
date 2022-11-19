package v1

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/neermitt/opsos/api/common"
	"github.com/neermitt/opsos/pkg/merge"
	"github.com/neermitt/opsos/pkg/utils"
	"github.com/pkg/errors"
)

type StacksSpec struct {
	BasePath      *string  `yaml:"base_path,omitempty" json:"base_path,omitempty" mapstructure:"base_path" validate:"required"`
	IncludedPaths []string `yaml:"included_paths,omitempty" json:"included_paths,omitempty" mapstructure:"included_paths" validate:"required"`
	ExcludedPaths []string `yaml:"excluded_paths,omitempty" json:"excluded_paths,omitempty" mapstructure:"excluded_paths"`
	NamePattern   *string  `yaml:"name_pattern,omitempty" json:"name_pattern,omitempty" mapstructure:"name_pattern" validate:"required"`
}

type WorkflowsSpec struct {
}

type ProviderSettings map[string]any

type ConfigSpec struct {
	BasePath  *string                     `yaml:"base_path,omitempty" json:"base_path,omitempty" mapstructure:"base_path" validate:"required"`
	Stacks    *StacksSpec                 `yaml:"stacks,omitempty" json:"stacks,omitempty" mapstructure:"stacks" validate:"required"`
	Workflows WorkflowsSpec               `yaml:"workflows,omitempty" json:"workflows,omitempty"`
	Providers map[string]ProviderSettings `yaml:",inline" json:",inline" mapstructure:",remain"`
}

type Config struct {
	common.Object `yaml:",inline" json:",inline"`
	Spec          ConfigSpec `yaml:"spec" json:"spec"`
}

func ReadAndMergeConfigsFromDirs(dirs []string) (*Config, error) {
	configs := make([]*Config, 0)
	for _, dir := range dirs {
		opsosConfigFileName := filepath.Join(dir, "opsos.yaml")
		if utils.FileExists(opsosConfigFileName) {
			config, err := ReadConfigFromFile(opsosConfigFileName)
			if err != nil {
				return nil, errors.Wrapf(err, "Invalid config file %s", opsosConfigFileName)
			}
			configs = append(configs, config)
		}
	}

	return MergeConfigs(configs)
}

func MergeConfigs(configs []*Config) (*Config, error) {
	switch len(configs) {
	case 0:
		return nil, nil
	case 1:
		return configs[0], nil
	}

	specs := make([]map[string]any, len(configs))
	for i, config := range configs {
		var err error
		specs[i], err = utils.ToMap(config.Spec)
		if err != nil {
			return nil, err
		}
	}
	mergedSpec, err := merge.Merge(specs)
	if err != nil {
		return nil, err
	}
	targetConfig := configs[len(configs)]
	err = utils.FromMap(mergedSpec, &targetConfig.Spec)
	if err != nil {
		return nil, err
	}
	return targetConfig, nil
}

func ReadConfigFromFile(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ReadConfig(file)
}

func ReadConfig(r io.Reader) (*Config, error) {
	var config Config

	err := utils.DecodeYaml(r, &config)
	if err != nil {
		return nil, err
	}
	err = validateConfig(config)
	if err != nil {
		return nil, err
	}
	return &config, err
}

func validateConfig(component Config) error {
	if component.ApiVersion != "opsos/v1" || component.Kind != "Configuration" {
		return fmt.Errorf("no resource found of type %s/%s", component.ApiVersion, component.Kind)
	}
	return nil
}
