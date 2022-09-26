package stack

import (
	"gopkg.in/yaml.v3"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/goburrow/cache"
	"github.com/mitchellh/mapstructure"
	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/merge"
	"github.com/neermitt/opsos/pkg/stack/schema"
	"github.com/neermitt/opsos/pkg/utils/fs"
	"github.com/spf13/afero"
)

type Stack struct {
	Name           string
	ComponentTypes map[string]ComponentMap
}

type ComponentMap struct {
	Components map[string]Component
}

type Component struct {
	Vars map[string]any
	Envs map[string]string
}

type StackProcessor interface {
	GetStackNames() ([]string, error)
	GetStack(name string) (*Stack, error)
	GetStacks(names []string) ([]*Stack, error)
}

func NewStackProcessor(source afero.Fs, includePaths []string, excludePaths []string) StackProcessor {
	sp := &stackProcessor{fs: source, fl: fs.NewMatcherFs(source, matcher(includePaths, excludePaths))}
	sp.cache = cache.NewLoadingCache(func(key cache.Key) (cache.Value, error) {
		return sp.loadAndProcessStackFile(key.(string))
	})
	return sp
}

func NewStackProcessorFromConfig(conf *config.Configuration) (StackProcessor, error) {
	stacksBasePath := path.Join(conf.BasePath, conf.Stacks.BasePath)
	stacksBaseAbsPath, err := filepath.Abs(stacksBasePath)
	if err != nil {
		return nil, nil
	}

	stackFS := afero.NewBasePathFs(afero.NewOsFs(), stacksBaseAbsPath)

	return NewStackProcessor(stackFS, conf.Stacks.IncludedPaths, conf.Stacks.ExcludedPaths), nil
}

type stackProcessor struct {
	fs    afero.Fs
	fl    afero.Fs
	cache cache.LoadingCache
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

func (sp *stackProcessor) GetStack(name string) (*Stack, error) {
	stackConfig, err := sp.loadAndProcessStackFile(name)
	if err != nil {
		return nil, err
	}
	return ProcessStackConfig(stackConfig)
}

func (sp *stackProcessor) GetStacks(names []string) ([]*Stack, error) {
	stkConfigs, err := sp.checkCacheOrLoadStackFiles(names)
	if err != nil {
		return nil, err
	}
	out := make([]*Stack, len(stkConfigs))
	for i, stk := range stkConfigs {
		stk, err := ProcessStackConfig(stk)
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

	importConfigs := make([]map[string]any, len(out.Import))

	imports, err := sp.checkCacheOrLoadStackFiles(out.Import)

	for i, imp := range imports {
		importConfigs[i] = imp.Config
	}

	out.Config, err = merge.Merge(append(importConfigs, out.Config))

	return out, nil
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

func ProcessStackConfig(stk *stack) (*Stack, error) {
	var stackConfig schema.StackConfig
	err := mapstructure.Decode(stk.Config, &stackConfig)
	if err != nil {
		return nil, err
	}
	componentTypeMap, err := processComponentType(stackConfig, "helmfile")
	if err != nil {
		return nil, err
	}
	return &Stack{Name: stk.name, ComponentTypes: map[string]ComponentMap{
		"helmfile": componentTypeMap,
	}}, nil
}

func processComponentType(stackConfig schema.StackConfig, componentType string) (ComponentMap, error) {
	componentTypeVars, err := merge.Merge([]map[string]any{stackConfig.Vars, stackConfig.ComponentTypeSettings[componentType].Vars})
	if err != nil {
		return ComponentMap{}, err
	}
	componentsMap := ComponentMap{Components: map[string]Component{}}
	for k, v := range stackConfig.Components.Types[componentType] {
		component, err := processComponentConfig(v, componentTypeVars)
		if err != nil {
			return ComponentMap{}, err
		}
		componentsMap.Components[k] = component
	}
	return componentsMap, nil
}

func processComponentConfig(conf schema.ComponentConfig, componentTypeVars map[string]any) (Component, error) {
	componentVars, err := merge.Merge([]map[string]any{componentTypeVars, conf.Vars})
	if err != nil {
		return Component{}, err
	}
	return Component{Vars: componentVars}, nil
}