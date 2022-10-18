package terraform

import (
	"fmt"
	"os"
	"strings"

	"github.com/neermitt/opsos/pkg/utils"
)

func ConstructWorkspaceName(execCtx ExecutionContext) (string, error) {
	var workspace string

	if execCtx.ComponentConfig.Metadata != nil && execCtx.ComponentConfig.Metadata.TerraformWorkspacePattern != nil {
		var err error
		workspace, err = utils.ProcessTemplate(*execCtx.ComponentConfig.Metadata.TerraformWorkspacePattern, execCtx.ComponentConfig.Vars)
		if err != nil {
			return "", err
		}
	} else if execCtx.ComponentConfig.Metadata != nil && execCtx.ComponentConfig.Metadata.TerraformWorkspace != nil {
		// Terraform workspace can be overridden per component in YAML config `metadata.terraform_workspace`
		workspace = *execCtx.ComponentConfig.Metadata.TerraformWorkspace
	} else {
		workspace = fmt.Sprintf("%s-%s", execCtx.Stack.Name, execCtx.ComponentName)
	}

	return strings.Replace(workspace, "/", "-", -1), nil
}

func SelectOrCreateWorkspace(execCtx ExecutionContext, workspace string) error {
	if err := SelectWorkspace(execCtx, workspace); err != nil {
		return CreateWorkspace(execCtx, workspace)
	}
	return nil
}

func SelectWorkspace(execCtx ExecutionContext, workspace string) error {

	args := []string{"workspace"}
	args = append(args, "select", workspace)

	cmdEnv, err := buildCommandEnvs(execCtx)
	if err != nil {
		return err
	}

	command := getCommand(execCtx)

	return utils.ExecuteShellCommand(execCtx.Context, command, args, utils.ExecOptions{
		DryRun:           execCtx.DryRun,
		Env:              cmdEnv,
		WorkingDirectory: execCtx.WorkingDir,
		StdOut:           os.Stdout,
	})
}

func CreateWorkspace(execCtx ExecutionContext, workspace string) error {

	args := []string{"workspace"}
	args = append(args, "new", workspace)

	cmdEnv, err := buildCommandEnvs(execCtx)
	if err != nil {
		return err
	}

	command := getCommand(execCtx)

	return utils.ExecuteShellCommand(execCtx.Context, command, args, utils.ExecOptions{
		DryRun:           execCtx.DryRun,
		Env:              cmdEnv,
		WorkingDirectory: execCtx.WorkingDir,
		StdOut:           os.Stdout,
	})
}
