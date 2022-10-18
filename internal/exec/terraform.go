package exec

import (
	"context"
	"errors"
	"github.com/neermitt/opsos/pkg/plugins/terraform"
	"github.com/neermitt/opsos/pkg/utils"
	"os"
)

type TerraformOptions struct {
	UsePlan     bool
	AutoApprove bool
	Destroy     bool
}

func ExecuteTerraformInit(ctx context.Context, stackName string, component string, additionalArgs []string, options TerraformOptions) error {
	exeCtx, err := terraform.NewExecutionContext(ctx, stackName, component, terraform.WithAdditionalArgs(additionalArgs))
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
	return terraformInit(exeCtx)
}

func terraformInit(exeCtx terraform.ExecutionContext) error {
	args := []string{"init"}
	if exeCtx.Config.Components.Terraform.InitRunReconfigure {
		args = append(args, "-reconfigure")
	}
	return terraform.ExecuteCommand(exeCtx, args)
}

func ExecuteTerraformPlan(ctx context.Context, stackName string, component string, additionalArgs []string, options TerraformOptions) error {
	exeCtx, err := terraform.NewExecutionContext(ctx, stackName, component, terraform.WithAdditionalArgs(additionalArgs))
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
	if err := terraformInit(exeCtx); err != nil {
		return err
	}

	workspace, err := terraform.ConstructWorkspaceName(exeCtx)
	if err != nil {
		return err
	}
	if err := terraform.SelectOrCreateWorkspace(exeCtx, workspace); err != nil {
		return err
	}

	args := []string{"plan"}
	args = append(args, "-var-file", exeCtx.VarFile, "-out", exeCtx.PlanFile)
	args = append(args, exeCtx.AdditionalArgs...)
	return terraform.ExecuteCommand(exeCtx, args)
}

func ExecuteTerraformApply(ctx context.Context, stackName string, component string, additionalArgs []string, options TerraformOptions) error {
	exeCtx, err := terraform.NewExecutionContext(ctx, stackName, component, terraform.WithAdditionalArgs(additionalArgs))
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
	if err := terraformInit(exeCtx); err != nil {
		return err
	}

	workspace, err := terraform.ConstructWorkspaceName(exeCtx)
	if err != nil {
		return err
	}
	if err := terraform.SelectOrCreateWorkspace(exeCtx, workspace); err != nil {
		return err
	}

	args := []string{"apply"}
	if options.UsePlan {
		args = append(args, exeCtx.PlanFile)
	} else {
		args = append(args, "-var-file", exeCtx.PlanFile)
	}

	if !utils.StringInSlice(terraform.AutoApproveFlag, exeCtx.AdditionalArgs) {
		if (exeCtx.Config.Components.Terraform.ApplyAutoApprove || options.AutoApprove) && !options.UsePlan {
			args = append(args, terraform.AutoApproveFlag)
		} else if os.Stdin == nil {
			return errors.New("'terraform apply' requires a user interaction, but it's running without `tty` or `stdin` attached.\nUse 'terraform apply -auto-approve' or 'terraform deploy' instead.")
		}
	}
	if options.Destroy {
		args = append(args, "-destroy")
	}
	args = append(args, exeCtx.AdditionalArgs...)

	if err := terraform.ExecuteCommand(exeCtx, args); err != nil {
		return err
	}

	return terraform.CleanPlanFile(exeCtx)
}

func ExecuteTerraformImport(ctx context.Context, stackName string, component string, additionalArgs []string, options TerraformOptions) error {
	exeCtx, err := terraform.NewExecutionContext(ctx, stackName, component, terraform.WithAdditionalArgs(additionalArgs))
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
	if err := terraformInit(exeCtx); err != nil {
		return err
	}

	workspace, err := terraform.ConstructWorkspaceName(exeCtx)
	if err != nil {
		return err
	}
	if err := terraform.SelectOrCreateWorkspace(exeCtx, workspace); err != nil {
		return err
	}

	args := []string{"import", "-var-file", exeCtx.VarFile}
	args = append(args, exeCtx.AdditionalArgs...)
	if err := terraform.ExecuteCommand(exeCtx, args); err != nil {
		return err
	}

	return terraform.CleanPlanFile(exeCtx)
}
