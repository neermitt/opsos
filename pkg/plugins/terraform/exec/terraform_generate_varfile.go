package exec

import (
	"context"

	"github.com/neermitt/opsos/pkg/plugins/terraform"
	"github.com/neermitt/opsos/pkg/stack"
)

type TerraformGenerateVarfileOptions struct {
	DryRun bool
	Format string
}

// ExecuteTerraformGenerateVarfile executes `terraform generate varfile` command
func ExecuteTerraformGenerateVarfile(ctx context.Context, stackName string, componentName string, options TerraformGenerateVarfileOptions) error {
	component := stack.Component{Type: terraform.ComponentType, Name: componentName}
	ctx = stack.SetStackName(ctx, stackName)
	ctx = stack.SetComponent(ctx, component)
	stk, err := stack.LoadStack(ctx, stack.LoadStackOptions{Stack: stackName, Component: &component})
	if err != nil {
		return err
	}
	ctx, err = terraform.NewExecutionContext(ctx, stk, component, true)
	if err != nil {
		return err
	}

	return terraform.GenerateVarFileFile(ctx, options.Format)
}
