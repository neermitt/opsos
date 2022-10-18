package terraform

import (
	"fmt"
	"strings"

	"github.com/neermitt/opsos/pkg/stack"
	"github.com/neermitt/opsos/pkg/utils"
)

func ConstructWorkspaceName(stack *stack.Stack, componentName string, config stack.ConfigWithMetadata) (string, error) {
	var workspace string

	if config.Metadata != nil && config.Metadata.TerraformWorkspacePattern != nil {
		var err error
		workspace, err = utils.ProcessTemplate(*config.Metadata.TerraformWorkspacePattern, config.Vars)
		if err != nil {
			return "", err
		}
	} else if config.Metadata != nil && config.Metadata.TerraformWorkspace != nil {
		// Terraform workspace can be overridden per component in YAML config `metadata.terraform_workspace`
		workspace = *config.Metadata.TerraformWorkspace
	} else {
		workspace = fmt.Sprintf("%s-%s", stack.Name, componentName)
	}

	return strings.Replace(workspace, "/", "-", -1), nil
}

func SelectOrCreateWorkspace(execCtx ExecutionContext) error {
	if err := ExecuteCommand(execCtx, []string{"workspace", "select", execCtx.WorkspaceName}); err != nil {
		return ExecuteCommand(execCtx, []string{"workspace", "new", execCtx.WorkspaceName})
	}
	return nil
}
