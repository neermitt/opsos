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
	Type                      string
	Component                 string
	Inherits                  []string
	TerraformWorkspace        string
	TerraformWorkspacePattern string
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

	// process Metadata Component Override
	if componentConfig.Metadata.Component != "" {
		componentConfig.Component = componentConfig.Metadata.Component
	}

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
	hierarchy, err := loadInheritanceTree(stackName, componentsConfigMap, componentName, true)
	if err != nil {
		return ConfigWithMetadata{}, err
	}

	hierarchy = unique(hierarchy)

	if len(hierarchy) != 0 {
		baseComponentConfigs := make([]ConfigWithMetadata, 0, len(hierarchy))
		for _, baseComponent := range hierarchy {
			if v, found := componentsConfigMap[baseComponent]; !found {
				return ConfigWithMetadata{}, fmt.Errorf("missing component %[2]s in stack %[1]s", stackName, baseComponent)
			} else {
				baseComponentConfigs = append(baseComponentConfigs, v)
			}
		}

		baseComponentsConfig, err := mergeConfigList(baseComponentConfigs)
		if err != nil {
			return ConfigWithMetadata{}, err
		}
		componentConfig = ConfigWithMetadata{Config: baseComponentsConfig, Metadata: componentConfig.Metadata}
	}

	// Update Component
	if componentConfig.Component == "" {
		componentConfig.Component = componentName
	}
	return componentConfig, nil
}

func loadInheritanceTree(stackName string, componentsConfigMap map[string]ConfigWithMetadata, componentName string, processInheritance bool) ([]string, error) {
	var componentConfig ConfigWithMetadata
	if v, found := componentsConfigMap[componentName]; !found {
		return nil, fmt.Errorf("missing component %[2]s in stack %[1]s", stackName, componentName)
	} else {
		componentConfig = v
	}
	componentHierarchy := make([]string, 0, 10)

	if componentConfig.Component != "" {
		baseComponentHierarchy, err := loadInheritanceTree(stackName, componentsConfigMap, componentConfig.Component, false)
		if err != nil {
			return nil, err
		}
		componentHierarchy = append(componentHierarchy, baseComponentHierarchy...)
	}

	if processInheritance {
		for _, baseComponent := range componentConfig.Metadata.Inherits {
			baseComponentHierarchy, err := loadInheritanceTree(stackName, componentsConfigMap, baseComponent, false)
			if err != nil {
				return nil, err
			}
			componentHierarchy = append(componentHierarchy, baseComponentHierarchy...)
		}
	}

	componentHierarchy = append(componentHierarchy, componentName)

	return componentHierarchy, nil
}

func mergeConfigList(configs []ConfigWithMetadata) (Config, error) {
	baseConfig := configs[0]

	for _, config := range configs[1:] {
		merged, err := mergeConfigs(baseConfig.Config, config.Config)
		if err != nil {
			return Config{}, err
		}
		baseConfig = ConfigWithMetadata{Config: merged}
	}

	return baseConfig.Config, nil
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

	if config1.Component != "" {
		mc["component"] = config1.Component
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

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
