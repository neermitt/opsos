package stack

import (
	"context"
	"errors"
	"fmt"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/utils"
)

type LoadStackOptions struct {
	Stack     string
	Component *Component
}

func LoadStack(ctx context.Context, options LoadStackOptions) (*Stack, error) {
	conf := config.GetConfig(ctx)
	stackProcessor, err := NewStackProcessorFromConfig(conf)
	if err != nil {
		return nil, nil
	}

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
	return stackProcessor.GetStack(options.Stack, options.Component)
}
