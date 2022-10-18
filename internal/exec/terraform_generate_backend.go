package exec

import (
	"context"

	"github.com/neermitt/opsos/pkg/plugins/terraform"
)

type TerraformGenerateBackendOptions struct {
	DryRun bool
	Format string
}

// ExecuteTerraformGenerateBackend executes `terraform generate backend` command
func ExecuteTerraformGenerateBackend(ctx context.Context, stackName string, component string, options TerraformGenerateBackendOptions) error {
	var opts []terraform.Option
	if options.DryRun {
		opts = append(opts, terraform.WithDryRun())
	}
	exeCtx, err := terraform.NewExecutionContext(ctx, stackName, component, opts...)
	if err != nil {
		return err
	}

	return terraform.GenerateBackendFile(exeCtx, options.Format)
}
