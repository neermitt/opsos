package terraform

import (
	"context"
	"fmt"
	"log"
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

func SelectOrCreateWorkspace(ctx context.Context) error {
	terraformSettings := GetTerraformSettings(ctx)
	log.Printf("[DEBUG] (terraform) Set workspacename %s", terraformSettings.WorkspaceName)
	if err := ExecuteCommand(ctx, []string{"workspace", "select", terraformSettings.WorkspaceName}); err != nil {
		return ExecuteCommand(ctx, []string{"workspace", "new", terraformSettings.WorkspaceName})
	}
	return nil
}
