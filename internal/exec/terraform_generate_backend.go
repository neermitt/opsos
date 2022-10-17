package exec

import (
	"context"
	"fmt"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/plugins/terraform"
	"github.com/neermitt/opsos/pkg/stack"
)

const (
	terraformComponentType = "terraform"
)

type TerraformGenerateBackendOptions struct {
	DryRun bool
	Format string
}

// ExecuteTerraformGenerateBackend executes `terraform generate backend` command
func ExecuteTerraformGenerateBackend(ctx context.Context, stackName string, component string, options TerraformGenerateBackendOptions) error {
	stk, err := stack.LoadStack(ctx, stack.LoadStackOptions{Stack: stackName, ComponentType: terraformComponentType, ComponentName: component})
	if err != nil {
		return err
	}

	terraformComponents, found := stk.Components[terraformComponentType]
	if !found {
		return fmt.Errorf("no terraform component found")
	}
	componentConfig, found := terraformComponents[component]
	if !found {
		return fmt.Errorf("terraform component %s not found", component)
	}

	conf := config.GetConfig(ctx)
	workingDir, _, err := getComponentWorkingDirectory(conf, terraformComponentType, componentConfig)
	if err != nil {
		return err
	}

	return terraform.GenerateBackendFile(terraform.ExecutionContext{
		Config:          config.GetConfig(ctx),
		Stack:           stk,
		ComponentName:   component,
		ComponentConfig: componentConfig,
		WorkingDir:      workingDir,
		DryRun:          options.DryRun,
	}, options.Format)
}
