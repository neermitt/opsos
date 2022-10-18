package exec

import (
	"context"
	"fmt"
	"os"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/plugins/terraform"
	"github.com/neermitt/opsos/pkg/stack"
	"github.com/neermitt/opsos/pkg/utils"
)

type TerraformCleanOptions struct {
	ClearDataDir bool
}

// ExecuteTerraformClean executes `terraform clean` command
func ExecuteTerraformClean(ctx context.Context, stackName string, component string, options TerraformCleanOptions) error {
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

	fmt.Print("Component backend config:\\n\\n")
	err = utils.GetFormatter("json")(os.Stdout, componentConfig.Vars)
	if err != nil {
		return err
	}

	conf := config.GetConfig(ctx)
	workingDir, _, err := getComponentWorkingDirectory(conf, terraformComponentType, componentConfig)
	if err != nil {
		return err
	}

	return terraform.Clean(terraform.ExecutionContext{
		Config:          config.GetConfig(ctx),
		Stack:           stk,
		ComponentName:   component,
		ComponentConfig: componentConfig,
		WorkingDir:      workingDir,
		DryRun:          false,
	}, options.ClearDataDir)
}
