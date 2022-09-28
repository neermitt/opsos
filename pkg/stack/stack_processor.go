package stack

import (
	"bytes"
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/goburrow/cache"
	"github.com/mitchellh/mapstructure"
	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/merge"
	"github.com/neermitt/opsos/pkg/plugins"
	"github.com/neermitt/opsos/pkg/stack/schema"
	"github.com/neermitt/opsos/pkg/utils/fs"
	"github.com/spf13/afero"
)

type Stack struct {
	Id                 string
	Name               string
	ComponentTypes     map[string]ComponentConfigMap
	Vars               map[string]any
	KubeConfigProvider string
}

type ComponentConfigMap map[string]any

type StackProcessor interface {
	GetStackNames(ctx context.Context) ([]string, error)
	GetStack(ctx context.Context, name string) (*Stack, error)
	GetStacks(ctx context.Context, names []string) ([]*Stack, error)
}

func NewStackProcessor(source afero.Fs, includePaths []string, excludePaths []string, stackNamePattern string) StackProcessor {
	tmpl := template.Must(template.New("stackNamePattern").Parse(stackNamePattern))

	sp := &stackProcessor{fs: source, fl: fs.NewMatcherFs(source, matcher(includePaths, excludePaths)), stackNameTemplate: tmpl}
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

	return NewStackProcessor(stackFS, conf.Stacks.IncludedPaths, conf.Stacks.ExcludedPaths, conf.Stacks.NamePattern), nil
}

type stackProcessor struct {
	fs                afero.Fs
	fl                afero.Fs
	cache             cache.LoadingCache
	stackNameTemplate *template.Template
}

func (sp *stackProcessor) GetStackNames(_ context.Context) ([]string, error) {
	files, err := fs.AllFiles(sp.fl)
	if err != nil {
		return nil, err
	}
	for i, val := range files {
		files[i] = strings.TrimSuffix(val, filepath.Ext(val))
	}
	return files, err
}

func (sp *stackProcessor) GetStack(ctx context.Context, name string) (*Stack, error) {
	stackConfig, err := sp.loadAndProcessStackFile(name)
	if err != nil {
		return nil, err
	}
	return sp.processStackConfig(ctx, stackConfig)
}

func (sp *stackProcessor) GetStacks(ctx context.Context, names []string) ([]*Stack, error) {
	stkConfigs, err := sp.checkCacheOrLoadStackFiles(names)
	if err != nil {
		return nil, err
	}
	out := make([]*Stack, len(stkConfigs))
	for i, stkConfig := range stkConfigs {
		stk, err := sp.processStackConfig(ctx, stkConfig)
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

func (sp *stackProcessor) processStackConfig(ctx context.Context, stk *stack) (*Stack, error) {
	var stackConfig schema.StackConfig
	err := mapstructure.Decode(stk.Config, &stackConfig)
	if err != nil {
		return nil, err
	}

	stackName, err := sp.getStackName(stackConfig.Vars)
	if err != nil {
		return nil, err
	}

	providers := plugins.GetProviders()

	componentTypes := make(map[string]ComponentConfigMap, len(providers))

	for _, providerName := range providers {
		componentTypeMap, err := sp.processComponentType(ctx, stackConfig, providerName)
		if err != nil {
			return nil, err
		}
		componentTypes[providerName] = componentTypeMap

	}

	return &Stack{Id: stk.name, Name: stackName, ComponentTypes: componentTypes, Vars: stackConfig.Vars, KubeConfigProvider: stackConfig.KubeConfigProvider}, nil
}

func (sp *stackProcessor) processComponentType(ctx context.Context, stackConfig schema.StackConfig, componentType string) (ComponentConfigMap, error) {

	provider, ok := plugins.GetProvider(ctx, componentType)
	if !ok {
		return ComponentConfigMap{}, fmt.Errorf("provider plugin not found for componentType: %s", componentType)
	}

	providerStackContext := provider.InitStackContext(context.Background(), stackConfig)

	componentsMap := ComponentConfigMap{}
	for k, v := range stackConfig.Components.Types[componentType] {
		componentExecConfig, err := provider.ProcessComponentConfig(providerStackContext, k, v)
		if err != nil {
			return ComponentConfigMap{}, err
		}
		componentsMap[k] = componentExecConfig
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
