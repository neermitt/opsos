package exec

import (
	"context"

	"github.com/neermitt/opsos/pkg/plugins/terraform"
)

type TerraformCleanOptions struct {
	ClearDataDir bool
}

// ExecuteTerraformClean executes `terraform clean` command
func ExecuteTerraformClean(ctx context.Context, stackName string, component string, options TerraformCleanOptions) error {
	exeCtx, err := terraform.NewExecutionContext(ctx, stackName, component)
	if err != nil {
		return err
	}
	return terraform.Clean(exeCtx, options.ClearDataDir)
}
