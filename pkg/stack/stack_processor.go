package stack

import (
	"path/filepath"
	"sync"

	"github.com/neermitt/opsos/pkg/merge"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

func ProcessYAMLFiles(afs afero.Fs, stackPaths []string) {

	count := len(stackPaths)

	var wg sync.WaitGroup

	wg.Add(count)

	for i, sp := range stackPaths {
		go func(i int, sp string) {
			defer wg.Done()

			ProcessYAMLFile(afs, sp)
		}(i, sp)
	}
}

type stack struct {
	Import []string       `yaml:"import,omitempty"`
	Config map[string]any `yaml:",inline"`
}

func ProcessYAMLFile(afs afero.Fs, sp string) (*stack, error) {
	filePath := sp
	if ext := filepath.Ext(sp); len(ext) == 0 {
		filePath = sp + ".yaml"
	}
	data, err := afero.ReadFile(afs, filePath)
	if err != nil {
		return nil, err
	}
	out := &stack{}
	err = yaml.Unmarshal(data, out)
	if err != nil {
		return nil, err
	}

	imports := make([]*stack, len(out.Import))
	importConfigs := make([]map[string]any, len(out.Import))

	for i, imp := range out.Import {
		is, err := ProcessYAMLFile(afs, imp)
		if err != nil {
			return nil, err
		}
		imports[i] = is
		importConfigs[i] = is.Config
	}

	out.Config, err = merge.Merge(append(importConfigs, out.Config))

	return out, nil
}
