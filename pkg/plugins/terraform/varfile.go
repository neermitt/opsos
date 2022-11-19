package terraform

import (
	"context"
	"log"
	"path"

	"github.com/neermitt/opsos/pkg/logging"
	"github.com/neermitt/opsos/pkg/stack"
	"github.com/neermitt/opsos/pkg/utils"
)

func GenerateVarFileFile(ctx context.Context, format string) error {
	terraformOptions := GetTerraformSettings(ctx)
	execOptions := utils.GetExecOptions(ctx)
	componentConfig := stack.GetComponentConfig(ctx)

	// Write varfile to file
	var varfilePath = path.Join(execOptions.WorkingDirectory, terraformOptions.VarFile)

	log.Printf("[INFO] (terraform) Writing the vars to file: %s", logging.Indent(varfilePath))

	if execOptions.DryRun {
		return nil
	}
	return utils.PrintOrWriteToFile(format, varfilePath, componentConfig.Vars, 0644)
}
