package terraform

import (
	"fmt"
	"os"
	"strings"

	"github.com/neermitt/opsos/pkg/utils"
)

func ExecutePlan(exeCtx ExecutionContext) error {
	args := []string{"plan"}
	planFile := constructPlanfileName(exeCtx)
	varFile := constructVarfileName(exeCtx, "")
	args = append(args, "-var-file", varFile, "-out", planFile)

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

// constructPlanfileName constructs the planfile name for a terraform component in a stack
func constructPlanfileName(exeCtx ExecutionContext) string {
	fmtdComponentFolderPrefix := strings.ReplaceAll(exeCtx.ComponentName, "/", "-")
	return fmt.Sprintf("%s-%s.planfile", exeCtx.Stack.Name, fmtdComponentFolderPrefix)
}
