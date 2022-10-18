package terraform

import (
	"os"

	"github.com/neermitt/opsos/pkg/utils"
)

func ExecuteImport(exeCtx ExecutionContext) error {
	args := []string{"import"}
	varFile := constructVarfileName(exeCtx, "")
	args = append(args, "-var-file", varFile)
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
