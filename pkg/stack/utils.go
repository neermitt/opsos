package stack

import (
	"context"
	"errors"
	"fmt"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/utils"
	"github.com/neermitt/opsos/pkg/utils/fs"
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

func LoadStack(ctx context.Context, options LoadStackOptions) (*Stack, error) {
	conf := config.GetConfig(ctx)
	stackProcessor, err := NewStackProcessorFromConfig(conf)
	if err != nil {
		return nil, nil
	}

	stackNames, err := stackProcessor.GetStackNames(ctx)
	if err != nil {
		return nil, nil
	}

	if options.Stack == "" {
		return nil, errors.New("stack must be specified")
	}

	if !utils.StringInSlice(options.Stack, stackNames) {
		return nil, fmt.Errorf("stack %s not found", options.Stack)
	}
	return stackProcessor.GetStack(ctx, options.Stack)
}
