package exec

import (
	"context"

	"github.com/neermitt/opsos/pkg/plugins/terraform"
)

type TerraformGenerateVarfileOptions struct {
	DryRun bool
	Format string
}

// ExecuteTerraformGenerateVarfile executes `terraform generate varfile` command
func ExecuteTerraformGenerateVarfile(ctx context.Context, stackName string, component string, options TerraformGenerateVarfileOptions) error {
	exeCtx, err := terraform.NewExecutionContext(ctx, stackName, component, options.DryRun)
	if err != nil {
		return err
	}

	return terraform.GenerateVarFileFile(exeCtx, options.Format)
}
