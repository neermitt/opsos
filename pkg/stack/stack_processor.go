package stack

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/goburrow/cache"
	"github.com/mitchellh/mapstructure"
	v1 "github.com/neermitt/opsos/api/v1"
	"github.com/neermitt/opsos/pkg/merge"
	"github.com/neermitt/opsos/pkg/stack/schema"
	"github.com/neermitt/opsos/pkg/utils/fs"
	"github.com/spf13/afero"
)

type Stack struct {
	Id         string
	Name       string
	Components map[string]ComponentConfigMap
	Vars       map[string]any
}

type ComponentConfigMap map[string]ConfigWithMetadata

type ConfigWithMetadata struct {
	Command                *string           `yaml:"command,omitempty" json:"command,omitempty" mapstructure:"command,omitempty"`
	Component              string            `yaml:"component,omitempty" json:"component,omitempty" mapstructure:"component,omitempty"`
	Vars                   map[string]any    `yaml:"vars,omitempty" json:"vars,omitempty"  mapstructure:"vars,omitempty"`
	Envs                   map[string]string `yaml:"env,omitempty" json:"env,omitempty"  mapstructure:"env,omitempty"`
	BackendType            *string           `yaml:"backend_type,omitempty" json:"backend_type,omitempty"  mapstructure:"backend_type,omitempty"`
	Backend                map[string]any    `yaml:"backend,omitempty" json:"backend,omitempty"  mapstructure:"backend,omitempty"`
	RemoteStateBackendType *string           `yaml:"remote_state_backend_type,omitempty" json:"remote_state_backend_type,omitempty"  mapstructure:"remote_state_backend_type,omitempty"`
	RemoteStateBackend     map[string]any    `yaml:"remote_state_backend,omitempty" json:"remote_state_backend,omitempty" mapstructure:"remote_state_backend,omitempty"`
	Settings               map[string]any    `yaml:"settings,omitempty" json:"settings,omitempty" mapstructure:"settings,omitempty"`
	Metadata               *schema.Metadata  `yaml:"metadata,omitempty" json:"metadata,omitempty" mapstructure:"metadata,omitempty"`
}

type GetStackOptions struct {
	ComponentTypes []string
	Components     []string
}

type StackProcessor interface {
	GetStackNames() ([]string, error)
	GetStack(name string, options GetStackOptions) (*Stack, error)
	GetStacks(names []string, options GetStackOptions) ([]*Stack, error)
}

func NewStackProcessor(source afero.Fs, includePaths []string, excludePaths []string, stackNamePattern string) StackProcessor {
	tmpl := template.Must(template.New("stackNamePattern").Parse(stackNamePattern))

	sp := &stackProcessor{fs: source, fl: fs.NewMatcherFs(source, fs.IncludeExcludeMatcher(includePaths, excludePaths)), stackNameTemplate: tmpl}
	sp.cache = cache.NewLoadingCache(func(key cache.Key) (cache.Value, error) {
		return sp.loadAndProcessStackFile(key.(string))
	})
	return sp
}

func NewStackProcessorFromConfig(conf *v1.ConfigSpec) (StackProcessor, error) {
	stacksBasePath := path.Join(*conf.BasePath, *conf.Stacks.BasePath)
	stacksBaseAbsPath, err := filepath.Abs(stacksBasePath)
	if err != nil {
		return nil, nil
	}

	stackFS := afero.NewBasePathFs(afero.NewOsFs(), stacksBaseAbsPath)

	return NewStackProcessor(stackFS, conf.Stacks.IncludedPaths, conf.Stacks.ExcludedPaths, *conf.Stacks.NamePattern), nil
}

type stackProcessor struct {
	fs                afero.Fs
	fl                afero.Fs
	cache             cache.LoadingCache
	stackNameTemplate *template.Template
}

func (sp *stackProcessor) GetStackNames() ([]string, error) {
	files, err := fs.AllFiles(sp.fl)
	if err != nil {
		return nil, err
	}
	for i, val := range files {
		files[i] = strings.TrimSuffix(val, filepath.Ext(val))
	}
	return files, err
}

func (sp *stackProcessor) GetStack(name string, options GetStackOptions) (*Stack, error) {
	stackConfig, err := sp.loadAndProcessStackFile(name)
	if err != nil {
		return nil, err
	}
	return sp.processStackConfig2(stackConfig, options)
}

func (sp *stackProcessor) GetStacks(names []string, options GetStackOptions) ([]*Stack, error) {
	stkConfigs, err := sp.checkCacheOrLoadStackFiles(names)
	if err != nil {
		return nil, err
	}
	out := make([]*Stack, len(stkConfigs))
	for i, stkConfig := range stkConfigs {
		stk, err := sp.processStackConfig2(stkConfig, options)
		if err != nil {
			return nil, err
		}
		out[i] = stk
	}
	return out, err
}

func (sp *stackProcessor) checkCacheOrLoadStackFiles(names []string) ([]*stack, error) {

	count := len(names)

	var wg sync.WaitGroup

	wg.Add(count)
	stacks := make([]*stack, count)

	var errorResult error

	for i, name := range names {
		go func(i int, name string) {
			defer wg.Done()

			stk, err := sp.checkCacheOrLoadStackFile(name)
			if err != nil {
				errorResult = err
				return
			}
			stacks[i] = stk
		}(i, name)
	}

	wg.Wait()

	if errorResult != nil {
		return nil, errorResult
	}

	return stacks, nil
}

func (sp *stackProcessor) checkCacheOrLoadStackFile(name string) (*stack, error) {
	val, err := sp.cache.Get(name)
	if err != nil {
		return nil, err
	}
	return val.(*stack), nil
}

func (sp *stackProcessor) loadAndProcessStackFile(name string) (*stack, error) {
	out, err := sp.loadStackFile(name)
	if err != nil {
		return nil, err
	}

	// Resolve stack files for imports
	importFiles, err := sp.resolveStackFiles(out.Import)
	if err != nil {
		return nil, err
	}

	importConfigs := make([]map[string]any, len(importFiles))

	imports, err := sp.checkCacheOrLoadStackFiles(importFiles)

	for i, imp := range imports {
		importConfigs[i] = imp.Config
	}

	out.Config, err = merge.Merge(append(importConfigs, out.Config))

	return out, nil
}

func (sp *stackProcessor) resolveStackFiles(filePatterns []string) ([]string, error) {
	matchedStackFiles := make([]string, 0, 2*len(filePatterns))
	for _, filePattern := range filePatterns {
		ext := filepath.Ext(filePattern)
		filePath := filePattern
		if ext := ext; len(ext) == 0 {
			filePath = filePattern + ".yaml"
		}
		match, err := afero.Glob(sp.fs, filePath)
		if err != nil {
			return nil, err
		}
		for _, m := range match {
			matchedStackFiles = append(matchedStackFiles, strings.TrimSuffix(m, ".yaml"))
		}
	}

	return matchedStackFiles, nil
}

func (sp *stackProcessor) loadStackFile(name string) (*stack, error) {
	filePath := name
	ext := filepath.Ext(name)
	if ext := ext; len(ext) == 0 {
		filePath = name + ".yaml"
	} else {
		name = strings.TrimSuffix(name, ext)
	}
	data, err := afero.ReadFile(sp.fs, filePath)
	if err != nil {
		return nil, err
	}

	out := &stack{name: name}
	err = yaml.Unmarshal(data, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type stack struct {
	name   string         `yaml:"_"`
	Import []string       `yaml:"import,omitempty"`
	Config map[string]any `yaml:",inline"`
}

func (sp *stackProcessor) processStackConfig(stk *stack, component *Component) (*Stack, error) {
	var stackConfig schema.StackConfig
	err := mapstructure.Decode(stk.Config, &stackConfig)
	if err != nil {
		return nil, err
	}

	stackName, err := sp.getStackName(stackConfig.Vars)
	if err != nil {
		return nil, err
	}

	var componentTypes []string
	if component != nil && component.Type != "" {
		componentTypes = []string{component.Type}
	} else {
		for k := range stackConfig.ComponentTypeSettings {
			componentTypes = append(componentTypes, k)
		}
	}

	processedComponentConfigs := make(map[string]ComponentConfigMap, len(componentTypes))

	for _, componentType := range componentTypes {
		componentTypeBaseConfig, err := getBaseConfigForComponentType(stackConfig, componentType)
		if err != nil {
			return nil, err
		}

		var componentsToProcess []string
		if component != nil && component.Name != "" {
			componentsToProcess = []string{component.Name}
		} else {
			for k := range stackConfig.Components.Types[componentType] {
				componentsToProcess = append(componentsToProcess, k)
			}
		}

		componentsMap := ComponentConfigMap{}
		for _, k := range componentsToProcess {
			componentProcessedConfig, err := processComponentConfigs(stk.name, componentTypeBaseConfig, stackConfig.Components.Types[componentType], k)
			if err != nil {
				return nil, err
			}
			configWithMetadata, err := toProcessedConfig(stk.name, k, componentProcessedConfig)
			if err != nil {
				return nil, err
			}
			componentsMap[k] = configWithMetadata
		}

		processedComponentConfigs[componentType] = componentsMap
	}

	return &Stack{Id: stk.name, Name: stackName, Components: processedComponentConfigs, Vars: stackConfig.Vars}, nil
}

func (sp *stackProcessor) processStackConfig2(stk *stack, options GetStackOptions) (*Stack, error) {
	var stackConfig schema.StackConfig
	err := mapstructure.Decode(stk.Config, &stackConfig)
	if err != nil {
		return nil, err
	}

	stackName, err := sp.getStackName(stackConfig.Vars)
	if err != nil {
		return nil, err
	}

	var componentTypes []string
	if len(options.ComponentTypes) != 0 {
		componentTypes = options.ComponentTypes
	} else {
		for k := range stackConfig.ComponentTypeSettings {
			componentTypes = append(componentTypes, k)
		}
	}

	processedComponentConfigs := make(map[string]ComponentConfigMap, len(componentTypes))

	for _, componentType := range componentTypes {
		componentTypeBaseConfig, err := getBaseConfigForComponentType(stackConfig, componentType)
		if err != nil {
			return nil, err
		}

		var componentsToProcess []string
		if len(options.Components) != 0 {
			componentsToProcess = options.Components
		} else {
			for k := range stackConfig.Components.Types[componentType] {
				componentsToProcess = append(componentsToProcess, k)
			}
		}

		componentsMap := ComponentConfigMap{}
		for _, k := range componentsToProcess {
			componentProcessedConfig, err := processComponentConfigs(stk.name, componentTypeBaseConfig, stackConfig.Components.Types[componentType], k)
			if err != nil {
				return nil, err
			}
			configWithMetadata, err := toProcessedConfig(stk.name, k, componentProcessedConfig)
			if err != nil {
				return nil, err
			}
			componentsMap[k] = configWithMetadata
		}

		processedComponentConfigs[componentType] = componentsMap
	}

	return &Stack{Id: stk.name, Name: stackName, Components: processedComponentConfigs, Vars: stackConfig.Vars}, nil
}

func (sp *stackProcessor) processComponentType(stackName string, stackConfig schema.StackConfig, componentType string) (ComponentConfigMap, error) {

	componentTypeBaseConfig, err := getBaseConfigForComponentType(stackConfig, componentType)
	if err != nil {
		return nil, err
	}

	componentsMap := ComponentConfigMap{}
	for k := range stackConfig.Components.Types[componentType] {
		componentProcessedConfig, err := processComponentConfigs(stackName, componentTypeBaseConfig, stackConfig.Components.Types[componentType], k)
		if err != nil {
			return ComponentConfigMap{}, err
		}
		processedConfig, err := toProcessedConfig(stackName, k, componentProcessedConfig)
		if err != nil {
			return ComponentConfigMap{}, err
		}
		componentsMap[k] = processedConfig
	}
	return componentsMap, nil
}

func (sp *stackProcessor) getStackName(vars map[string]any) (string, error) {
	var buff bytes.Buffer
	if err := sp.stackNameTemplate.Execute(&buff, vars); err != nil {
		return "", err
	}
	return buff.String(), nil
}

func getBaseConfigForComponentType(config schema.StackConfig, componentType string) (schema.Config, error) {
	globalConfig := schema.Config{
		Vars:     config.Vars,
		Envs:     config.Envs,
		Settings: config.Settings,
	}
	componentTypeSettings := config.ComponentTypeSettings[componentType]
	backend := componentTypeSettings.BackendType
	remoteStateBackend := componentTypeSettings.RemoteStateBackendType
	stackComponentConfig := schema.Config{
		Vars:                      componentTypeSettings.Vars,
		Envs:                      componentTypeSettings.Envs,
		BackendType:               &backend,
		BackendConfigs:            componentTypeSettings.Backend,
		RemoteStateBackendType:    &remoteStateBackend,
		RemoteStateBackendConfigs: componentTypeSettings.RemoteStateBackend,
		Settings:                  componentTypeSettings.Settings,
	}
	return mergeConfigs(globalConfig, stackComponentConfig)
}

func toProcessedConfig(stackName string, componentName string, componentProcessedConfig *schema.ConfigWithMetadata) (ConfigWithMetadata, error) {
	var processedBackend, processedRemoteStateBackend map[string]any
	if componentProcessedConfig.BackendType != nil && *componentProcessedConfig.BackendType != "" {
		var found bool
		processedBackend, found = componentProcessedConfig.BackendConfigs[*componentProcessedConfig.BackendType].(map[string]any)
		if !found {
			return ConfigWithMetadata{}, fmt.Errorf("backend %[3]s config not found for component %[2]s in stack %[1]s", stackName, componentName, *componentProcessedConfig.BackendType)
		}
		processedRemoteStateBackend, found = componentProcessedConfig.RemoteStateBackendConfigs[*componentProcessedConfig.RemoteStateBackendType].(map[string]any)
		if !found {
			return ConfigWithMetadata{}, fmt.Errorf("remote_state_backend %[3]s config not found for component %[2]s in stack %[1]s", stackName, componentName, *componentProcessedConfig.RemoteStateBackendType)
		}
	}
	configWithMetadata := ConfigWithMetadata{
		Command:                componentProcessedConfig.Command,
		Component:              *componentProcessedConfig.Component,
		Vars:                   componentProcessedConfig.Vars,
		Envs:                   componentProcessedConfig.Envs,
		BackendType:            componentProcessedConfig.BackendType,
		Backend:                processedBackend,
		RemoteStateBackendType: componentProcessedConfig.RemoteStateBackendType,
		RemoteStateBackend:     processedRemoteStateBackend,
		Settings:               componentProcessedConfig.Settings,
		Metadata:               componentProcessedConfig.Metadata,
	}
	return configWithMetadata, nil
}
