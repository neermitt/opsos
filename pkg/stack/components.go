package stack

import (
	"fmt"
	"github.com/neermitt/opsos/pkg/merge"
	"github.com/neermitt/opsos/pkg/stack/schema"
	"github.com/neermitt/opsos/pkg/utils"
	"github.com/pkg/errors"
)

func processComponentConfigs(stackName string, baseConfig schema.Config, componentsConfigMap map[string]schema.ConfigWithMetadata, componentName string) (*schema.ConfigWithMetadata, error) {
	componentConfig, err := loadComponentConfig(stackName, componentsConfigMap, componentName)
	if err != nil {
		return nil, err
	}

	// merge with base config
	mc, err := mergeConfigs(baseConfig, componentConfig.Config)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to merge config for %[2]s in stack %[1]s", stackName, componentName))
	}
	componentConfig = schema.ConfigWithMetadata{Config: mc, Metadata: componentConfig.Metadata}

	// process Metadata Component Override
	if componentConfig.Metadata != nil && componentConfig.Metadata.Component != nil {
		componentConfig.Component = componentConfig.Metadata.Component
	}

	// process remoteBackend
	if componentConfig.RemoteStateBackendType == nil || *componentConfig.RemoteStateBackendType == "" {
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

func loadComponentConfig(stackName string, componentsConfigMap map[string]schema.ConfigWithMetadata, componentName string) (schema.ConfigWithMetadata, error) {
	var componentConfig schema.ConfigWithMetadata
	if v, found := componentsConfigMap[componentName]; !found {
		return schema.ConfigWithMetadata{}, fmt.Errorf("missing component %[2]s in stack %[1]s", stackName, componentName)
	} else {
		componentConfig = v
	}

	// check inheritance
	hierarchy, err := loadInheritanceTree(stackName, componentsConfigMap, componentName, true)
	if err != nil {
		return schema.ConfigWithMetadata{}, err
	}

	hierarchy = utils.Unique(hierarchy)

	if len(hierarchy) != 0 {
		baseComponentConfigs := make([]schema.ConfigWithMetadata, 0, len(hierarchy))
		for _, baseComponent := range hierarchy {
			if v, found := componentsConfigMap[baseComponent]; !found {
				return schema.ConfigWithMetadata{}, fmt.Errorf("missing component %[2]s in stack %[1]s", stackName, baseComponent)
			} else {
				baseComponentConfigs = append(baseComponentConfigs, v)
			}
		}

		baseComponentsConfig, err := mergeConfigList(baseComponentConfigs)
		if err != nil {
			return schema.ConfigWithMetadata{}, err
		}
		componentConfig = schema.ConfigWithMetadata{Config: baseComponentsConfig, Metadata: componentConfig.Metadata}
	}

	// Update Component
	if componentConfig.Component == nil {
		componentConfig.Component = &componentName
	}
	return componentConfig, nil
}

func loadInheritanceTree(stackName string, componentsConfigMap map[string]schema.ConfigWithMetadata, componentName string, processInheritance bool) ([]string, error) {
	var componentConfig schema.ConfigWithMetadata
	if v, found := componentsConfigMap[componentName]; !found {
		return nil, fmt.Errorf("missing component %[2]s in stack %[1]s", stackName, componentName)
	} else {
		componentConfig = v
	}
	componentHierarchy := make([]string, 0, 10)

	if componentConfig.Component != nil {
		baseComponentHierarchy, err := loadInheritanceTree(stackName, componentsConfigMap, *componentConfig.Component, false)
		if err != nil {
			return nil, err
		}
		componentHierarchy = append(componentHierarchy, baseComponentHierarchy...)
	}

	if processInheritance && componentConfig.Metadata != nil {
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

func mergeConfigList(configs []schema.ConfigWithMetadata) (schema.Config, error) {
	baseConfig := configs[0]

	for _, conf := range configs[1:] {
		merged, err := mergeConfigs(baseConfig.Config, conf.Config)
		if err != nil {
			return schema.Config{}, err
		}
		baseConfig = schema.ConfigWithMetadata{Config: merged}
	}

	return baseConfig.Config, nil
}

func mergeConfigs(config1 schema.Config, config2 schema.Config) (schema.Config, error) {
	c1, err := config1.ToMap()
	if err != nil {
		return schema.Config{}, err
	}
	c2, err := config2.ToMap()
	if err != nil {
		return schema.Config{}, err
	}
	mc, err := merge.Merge([]map[string]any{c1, c2})
	if err != nil {
		return schema.Config{}, err
	}

	if config1.Component != nil {
		mc["component"] = config1.Component
	}

	return schema.NewConfigFromMap(mc)
}
