package exec

import (
	"context"
	"errors"
	"os"

	"github.com/neermitt/opsos/pkg/plugins/terraform"
	"github.com/neermitt/opsos/pkg/utils"
)

type TerraformOptions struct {
	UsePlan                   bool
	AutoApprove               bool
	Destroy                   bool
	SkipInit                  bool
	CleanPlanFileOnCompletion bool
}

func ExecuteTerraformInit(ctx context.Context, stackName string, component string, additionalArgs []string, options TerraformOptions) error {
	return executeTerraform(ctx, stackName, component, additionalArgs, options,
		func(exeCtx terraform.ExecutionContext, additionalArgs []string, options TerraformOptions) ([]string, error) {
			return prepareInitCommandArgs(exeCtx, additionalArgs), nil
		},
	)
}

func ExecuteTerraformPlan(ctx context.Context, stackName string, component string, additionalArgs []string, options TerraformOptions) error {
	return executeTerraform(ctx, stackName, component, additionalArgs, options,
		func(exeCtx terraform.ExecutionContext, additionalArgs []string, options TerraformOptions) ([]string, error) {
			args := []string{"plan", "-var-file", exeCtx.VarFile, "-out", exeCtx.PlanFile}
			args = append(args, additionalArgs...)
			return args, nil
		},
	)
}

func ExecuteTerraformApply(ctx context.Context, stackName string, component string, additionalArgs []string, options TerraformOptions) error {
	return executeTerraform(ctx, stackName, component, additionalArgs, options,
		func(exeCtx terraform.ExecutionContext, additionalArgs []string, options TerraformOptions) ([]string, error) {
			args := []string{"apply"}
			if options.UsePlan {
				args = append(args, exeCtx.PlanFile)
			} else {
				args = append(args, "-var-file", exeCtx.PlanFile)
			}

			if !utils.StringInSlice(terraform.AutoApproveFlag, additionalArgs) {
				if (exeCtx.Config.Components.Terraform.ApplyAutoApprove || options.AutoApprove) && !options.UsePlan {
					args = append(args, terraform.AutoApproveFlag)
				} else if os.Stdin == nil {
					return nil, errors.New("'terraform apply' requires a user interaction, but it's running without `tty` or `stdin` attached.\nUse 'terraform apply -auto-approve' or 'terraform deploy' instead.")
				}
			}
			if options.Destroy {
				args = append(args, "-destroy")
			}
			args = append(args, additionalArgs...)
			return args, nil
		},
	)
}

func ExecuteTerraformImport(ctx context.Context, stackName string, component string, additionalArgs []string, options TerraformOptions) error {
	return executeTerraform(ctx, stackName, component, additionalArgs, options,
		func(exeCtx terraform.ExecutionContext, additionalArgs []string, options TerraformOptions) ([]string, error) {
			args := []string{"import", "-var-file", exeCtx.VarFile}
			args = append(args, additionalArgs...)
			return args, nil
		},
	)
}

type prepareCommandArgs func(exeCtx terraform.ExecutionContext, additionalArgs []string, options TerraformOptions) ([]string, error)

func executeTerraform(ctx context.Context, stackName string, component string, additionalArgs []string, options TerraformOptions, prepArgs prepareCommandArgs) error {
	exeCtx, err := terraform.NewExecutionContext(ctx, stackName, component)
	if err != nil {
		return err
	}
	if err := terraform.GenerateVarFileFile(exeCtx, "json"); err != nil {
		return err
	}
	if exeCtx.Config.Components.Terraform.AutoGenerateBackendFile {
		if err := terraform.GenerateBackendFile(exeCtx, "json"); err != nil {
			return err
		}
	}
	if !options.SkipInit {
		initArgs := prepareInitCommandArgs(exeCtx, nil)
		if err := terraform.ExecuteCommand(exeCtx, initArgs); err != nil {
			return err
		}
		if err := terraform.SelectOrCreateWorkspace(exeCtx); err != nil {
			return err
		}
	}

	args, err := prepArgs(exeCtx, additionalArgs, options)
	if err != nil {
		return err
	}

	if err := terraform.ExecuteCommand(exeCtx, args); err != nil {
		return err
	}
	if options.CleanPlanFileOnCompletion {
		return terraform.CleanPlanFile(exeCtx)
	}
	return nil
}

func prepareInitCommandArgs(exeCtx terraform.ExecutionContext, additionalArgs []string) []string {
	args := []string{"init"}
	if exeCtx.Config.Components.Terraform.InitRunReconfigure {
		args = append(args, "-reconfigure")
	}
	args = append(args, additionalArgs...)
	return args
}
