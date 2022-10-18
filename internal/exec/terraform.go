package exec

import (
	"context"
	"errors"
	"os"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/plugins/terraform"
	"github.com/neermitt/opsos/pkg/utils"
)

type TerraformOptions struct {
	Command                   string
	RequiresVarFile           bool
	UsePlan                   bool
	DryRun                    bool
	AutoApprove               bool
	Destroy                   bool
	SkipInit                  bool
	SkipWorkspace             bool
	CleanPlanFileOnCompletion bool
}

func ExecuteTerraform(ctx context.Context, stackName string, component string, additionalArgs []string, options TerraformOptions) error {
	exeCtx, err := terraform.NewExecutionContext(ctx, stackName, component, options.DryRun)
	if err != nil {
		return err
	}

	if options.RequiresVarFile {
		if err := terraform.GenerateVarFileFile(exeCtx, "json"); err != nil {
			return err
		}
	}

	conf := config.GetConfig(ctx)
	if conf.Components.Terraform.AutoGenerateBackendFile {
		if err := terraform.GenerateBackendFile(exeCtx, "json"); err != nil {
			return err
		}
	}

	if !options.SkipInit {
		initArgs, err := prepareArgs(exeCtx, nil, TerraformOptions{Command: "init"})
		if err != nil {
			return err
		}
		if err := terraform.ExecuteCommand(exeCtx, initArgs); err != nil {
			return err
		}
	}
	if !options.SkipWorkspace {
		if err := terraform.SelectOrCreateWorkspace(exeCtx); err != nil {
			return err
		}
	}

	args, err := prepareArgs(exeCtx, additionalArgs, options)
	if err != nil {
		return err
	}

	if options.Command == "shell" {
		return terraform.ExecuteShell(exeCtx)
	}
	if err := terraform.ExecuteCommand(exeCtx, args); err != nil {
		return err
	}
	if options.CleanPlanFileOnCompletion {
		return terraform.CleanPlanFile(exeCtx)
	}
	return nil
}

func prepareArgs(exeCtx terraform.ExecutionContext, additionalArgs []string, options TerraformOptions) ([]string, error) {
	conf := config.GetConfig(exeCtx.Context)
	args := []string{options.Command}

	if options.UsePlan {
		args = append(args, exeCtx.PlanFile)
	} else if options.RequiresVarFile {
		args = append(args, "-var-file", exeCtx.VarFile)
	}

	switch options.Command {
	case "plan":
		args = append(args, "-out", exeCtx.PlanFile)
	case "init":
		if conf.Components.Terraform.InitRunReconfigure {
			args = append(args, "-reconfigure")
		}
	case "apply":
		if !utils.StringInSlice(terraform.AutoApproveFlag, additionalArgs) {
			if (conf.Components.Terraform.ApplyAutoApprove || options.AutoApprove) && !options.UsePlan {
				args = append(args, terraform.AutoApproveFlag)
			} else if os.Stdin == nil {
				return nil, errors.New("'terraform apply' requires a user interaction, but it's running without `tty` or `stdin` attached.\nUse 'terraform apply -auto-approve' or 'terraform deploy' instead.")
			}
		}
		if options.Destroy {
			args = append(args, "-destroy")
		}
	}
	args = append(args, additionalArgs...)
	return args, nil
}
