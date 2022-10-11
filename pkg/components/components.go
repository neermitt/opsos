package components

import (
	"fmt"
	"gopkg.in/yaml.v3"

	"github.com/neermitt/opsos/pkg/merge"
	"github.com/pkg/errors"
)

type Config struct {
	Component                 string `yaml:"component,omitempty"`
	Vars                      map[string]any
	Envs                      map[string]string
	BackendType               string `yaml:"backend_type,omitempty"`
	BackendConfigs            map[string]any
	RemoteStateBackendType    string `yaml:"remote_state_backend_type,omitempty"`
	RemoteStateBackendConfigs map[string]any
	Settings                  map[string]any
}

type ConfigWithMetadata struct {
	Config   `yaml:",inline"`
	Metadata Metadata
}

type Metadata struct {
	Type      string
	Component string
	Inherits  []string
}

func processComponentConfigs(stackName string, baseConfig Config, componentsConfigMap map[string]ConfigWithMetadata, componentName string) (*ConfigWithMetadata, error) {
	componentConfig, err := loadComponentConfig(stackName, componentsConfigMap, componentName)
	if err != nil {
		return nil, err
	}

	// merge with base config
	mc, err := mergeConfigs(baseConfig, componentConfig.Config)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to merge config for %[2]s in stack %[1]s", stackName, componentName))
	}
	componentConfig = ConfigWithMetadata{Config: mc, Metadata: componentConfig.Metadata}

	// process remoteBackend
	if componentConfig.RemoteStateBackendType == "" {
		componentConfig.RemoteStateBackendType = componentConfig.BackendType
	}

	if componentConfig.RemoteStateBackendConfigs == nil {
		componentConfig.RemoteStateBackendConfigs = componentConfig.BackendConfigs
	} else {
		mergedConfig, err := merge.Merge([]map[string]any{componentConfig.BackendConfigs, componentConfig.RemoteStateBackendConfigs})
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to merge remote_state_backend and backend configs for %[2]s in stack %[1]s", stackName, componentName))
		}
		componentConfig.RemoteStateBackendConfigs = mergedConfig
	}

	return &componentConfig, nil
}

func loadComponentConfig(stackName string, componentsConfigMap map[string]ConfigWithMetadata, componentName string) (ConfigWithMetadata, error) {
	var componentConfig ConfigWithMetadata
	if v, found := componentsConfigMap[componentName]; !found {
		return ConfigWithMetadata{}, fmt.Errorf("missing component %[2]s in stack %[1]s", stackName, componentName)
	} else {
		componentConfig = v
	}

	// check inheritance
	if componentConfig.Component != "" {
		baseComponentConfig, err := loadComponentConfig(stackName, componentsConfigMap, componentConfig.Component)
		if err != nil {
			return ConfigWithMetadata{}, fmt.Errorf("missing component %[3]s in stack %[1]s while inheriting from %[2]s", stackName, componentName, componentConfig.Component)
		}

		mc, err := mergeConfigs(baseComponentConfig.Config, componentConfig.Config)
		if err != nil {
			return ConfigWithMetadata{}, errors.Wrap(err, fmt.Sprintf("failed to merge config for %[3]s and %[2]s in stack %[1]s", stackName, componentName, componentConfig.Component))
		}
		mc.Component = baseComponentConfig.Component
		componentConfig = ConfigWithMetadata{Config: mc, Metadata: componentConfig.Metadata}
	}

	// Update Component
	if componentConfig.Component == "" {
		componentConfig.Component = componentName
	}
	return componentConfig, nil
}

func mergeConfigs(config1 Config, config2 Config) (Config, error) {
	c1, err := toMap(config1)
	if err != nil {
		return Config{}, err
	}
	c2, err := toMap(config2)
	if err != nil {
		return Config{}, err
	}
	mc, err := merge.Merge([]map[string]any{c1, c2})
	if err != nil {
		return Config{}, err
	}

	if config2.BackendType == "" {
		mc["backend_type"] = config1.BackendType
	}
	if config2.RemoteStateBackendType == "" {
		mc["remote_state_backend_type"] = config1.RemoteStateBackendType
	}

	return fromMap(mc)
}
func toMap(config Config) (map[string]any, error) {
	yamlCurrent, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	var dataCurrent map[string]any
	if err = yaml.Unmarshal(yamlCurrent, &dataCurrent); err != nil {
		return nil, err
	}

	return dataCurrent, nil
}

func fromMap(config map[string]any) (Config, error) {
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
