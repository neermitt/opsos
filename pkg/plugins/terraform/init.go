package terraform

import (
	"os"

	"github.com/neermitt/opsos/pkg/utils"
)

type InitOptions struct {
	Reconfigure bool
}

func ExecuteInit(exeCtx ExecutionContext, options InitOptions) error {
	args := []string{"init"}
	if options.Reconfigure {
		args = append(args, "-reconfigure")
	}

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
