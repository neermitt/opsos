package terraform

import (
	"fmt"
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
	if err := ExecuteCommand(execCtx, []string{"workspace", "select", workspace}); err != nil {
		return ExecuteCommand(execCtx, []string{"workspace", "new", workspace})
	}
	return nil
}
