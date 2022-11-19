package exec

import (
	"context"

	"github.com/neermitt/opsos/pkg/components"
	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/plugins/terraform"
	"github.com/neermitt/opsos/pkg/stack"
	"github.com/neermitt/opsos/pkg/utils"
)

type ComponentInitOptions struct {
	ComponentType string
	ComponentName string
	DryRun        bool
}

func ExecuteComponentInit(ctx context.Context, options ComponentInitOptions) error {
	conf := config.GetConfig(ctx)
	workingDir := components.GetWorkingDirectory(conf, options.ComponentType, options.ComponentName)

	return components.PrepareComponent(ctx, workingDir, workingDir, components.PrepareComponentOptions{DryRun: options.DryRun})
}

func ExecuteStackComponentsInit(ctx context.Context, stackName string, options ComponentInitOptions) error {
	component := stack.Component{Type: options.ComponentType, Name: options.ComponentName}
	ctx = stack.SetStackName(ctx, stackName)
	stk, err := stack.LoadStack(ctx, stack.LoadStackOptions{Stack: stackName, Component: &component})
	if err != nil {
		return err
	}

	return initStackComponents(ctx, stk, options)

}

func initStackComponents(ctx context.Context, stk *stack.Stack, options ComponentInitOptions) error {
	for componentType, componentMap := range stk.Components {
		for componentName := range componentMap {
			componentRef := stack.Component{Type: componentType, Name: componentName}
			ctx, err := terraform.NewExecutionContext(ctx, stk, componentRef, options.DryRun)
			if err != nil {
				return err
			}

			execOptions := utils.GetExecOptions(ctx)
			err = components.PrepareComponent(ctx, execOptions.WorkingDirectory, execOptions.WorkingDirectory,
				components.PrepareComponentOptions{DryRun: options.DryRun})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
