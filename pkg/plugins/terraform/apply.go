package terraform

import (
	"os"

	"github.com/neermitt/opsos/pkg/utils"
	"github.com/pkg/errors"
)

const autoApproveFlag = "-auto-approve"

type ApplyOptions struct {
	UsePlan     bool
	AutoApprove bool
	Destroy     bool
}

func ExecuteApply(exeCtx ExecutionContext, options ApplyOptions) error {
	args := []string{"apply"}
	if options.UsePlan {
		planFile := constructPlanfileName(exeCtx)
		args = append(args, planFile)
	} else {
		varFile := constructVarfileName(exeCtx, "")
		args = append(args, "-var-file", varFile)
	}

	if !utils.StringInSlice(autoApproveFlag, exeCtx.AdditionalArgs) {
		if (exeCtx.Config.Components.Terraform.ApplyAutoApprove || options.AutoApprove) && !options.UsePlan {
			args = append(args, autoApproveFlag)
		} else if os.Stdin == nil {
			return errors.New("'terraform apply' requires a user interaction, but it's running without `tty` or `stdin` attached.\nUse 'terraform apply -auto-approve' or 'terraform deploy' instead.")
		}
	}
	if options.Destroy {
		args = append(args, "-destroy")
	}

	args = append(args, exeCtx.AdditionalArgs...)

	cmdEnv, err := buildCommandEnvs(exeCtx)
	if err != nil {
		return err
	}

	command := getCommand(exeCtx)

	return utils.ExecuteShellCommand(exeCtx.Context, command, args, utils.ExecOptions{
		DryRun:           exeCtx.DryRun,
		Env:              cmdEnv,
		WorkingDirectory: exeCtx.WorkingDir,
		StdOut:           os.Stdout,
	})
}
