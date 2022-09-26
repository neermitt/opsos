package stack

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/utils"
	"github.com/neermitt/opsos/pkg/utils/fs"
	"github.com/spf13/afero"
)

func matcher(includedPaths []string, excludedPaths []string) fs.Matcher {
	includeMatcher := globMatchers(includedPaths)
	excludedMatcher := globMatchers(excludedPaths)

	return fs.And(includeMatcher, fs.Not(excludedMatcher))
}

func globMatchers(paths []string) fs.Matcher {
	matchers := make([]fs.Matcher, 0)
	for _, p := range paths {
		matchers = append(matchers, fs.Glob(p))
	}
	return fs.Or(matchers...)
}

type LoadStackOptions struct {
	Stack string
}

func LoadStack(conf *config.Configuration, options LoadStackOptions) (*Stack, error) {
	stacksBasePath := path.Join(conf.BasePath, conf.Stacks.BasePath)
	stacksBaseAbsPath, err := filepath.Abs(stacksBasePath)
	if err != nil {
		return nil, nil
	}

	stackFS := afero.NewBasePathFs(afero.NewOsFs(), stacksBaseAbsPath)

	stackProcessor := NewStackProcessor(stackFS, conf.Stacks.IncludedPaths, conf.Stacks.ExcludedPaths)
	stackNames, err := stackProcessor.GetStackNames()
	if err != nil {
		return nil, nil
	}

	if options.Stack == "" {
		return nil, errors.New("stack must be specified")
	}

	if !utils.StringInSlice(options.Stack, stackNames) {
		return nil, fmt.Errorf("stack %s not found", options.Stack)
	}
	return stackProcessor.GetStack(options.Stack)
}
