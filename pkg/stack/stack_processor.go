package stack

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"sync"

	"github.com/goburrow/cache"
	"github.com/neermitt/opsos/pkg/merge"
	"github.com/spf13/afero"
)

type Stack interface {
	Name() string
}

type StackProcessor interface {
	GetStack(name string) (Stack, error)
	GetStacks(names []string) ([]Stack, error)
}

func NewStackProcessor(fs afero.Fs) StackProcessor {
	sp := &stackProcessor{fs: fs}
	sp.cache = cache.NewLoadingCache(func(key cache.Key) (cache.Value, error) {
		fmt.Println(key)
		return sp.loadAndProcessStackFile(key.(string))
	}, cache.WithRemovalListener(func(key cache.Key, value cache.Value) {
		fmt.Printf("Remove: %s \n", key)
	}))
	return sp
}

type stackProcessor struct {
	fs    afero.Fs
	cache cache.LoadingCache
}

func (sp *stackProcessor) GetStack(name string) (Stack, error) {
	return sp.loadAndProcessStackFile(name)
}

func (sp *stackProcessor) GetStacks(names []string) ([]Stack, error) {
	stks, err := sp.checkCacheOrLoadStackFiles(names)
	if err != nil {
		return nil, err
	}
	out := make([]Stack, len(stks))
	for i, stk := range stks {
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
	if ext := filepath.Ext(name); len(ext) == 0 {
		filePath = name + ".yaml"
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

func (s *stack) Name() string {
	return s.name
}
