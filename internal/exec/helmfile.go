package exec

import (
	"context"
	"strings"

	"github.com/neermitt/opsos/pkg/plugins/helmfile"
	"github.com/neermitt/opsos/pkg/stack"
)

type HelmfileExecOptions struct {
	DryRun     bool
	GlobalArgs string
}

func ExecHelmfile(ctx context.Context, command string, stackName string, componentName string, additionalArgs []string, options HelmfileExecOptions) error {

	component := stack.Component{Type: helmfile.ComponentType, Name: componentName}
	ctx = stack.SetStackName(ctx, stackName)
	ctx = stack.SetComponent(ctx, component)

	stk, err := stack.LoadStack(ctx, stack.LoadStackOptions{Stack: stackName, Component: &component})
	if err != nil {
		return err
	}

	globalArgs := strings.Fields(options.GlobalArgs)
	return helmfile.ExecHelmfileCommand(ctx, command, stk, globalArgs, additionalArgs, options.DryRun)
}
