package exec

import (
	"context"

	"github.com/neermitt/opsos/pkg/components"
	"github.com/neermitt/opsos/pkg/plugins/terraform"
	"github.com/neermitt/opsos/pkg/stack"
	"github.com/neermitt/opsos/pkg/utils"
)

type ComponentPullOptions struct {
	ComponentType string
	DryRun        bool
}

func ExecuteComponentPull(ctx context.Context, stackName string, componentName string, options ComponentPullOptions) error {
	component := stack.Component{Type: options.ComponentType, Name: componentName}
	ctx = stack.SetStackName(ctx, stackName)
	ctx = stack.SetComponent(ctx, component)
	stk, err := stack.LoadStack(ctx, stack.LoadStackOptions{Stack: stackName, Component: &component})
	if err != nil {
		return err
	}
	ctx, err = terraform.NewExecutionContext(ctx, stk, component, options.DryRun)
	if err != nil {
		return err
	}

	execOptions := utils.GetExecOptions(ctx)
	return components.PrepareComponent(ctx, execOptions.WorkingDirectory, execOptions.WorkingDirectory,
		components.PrepareComponentOptions{DryRun: options.DryRun})
}
