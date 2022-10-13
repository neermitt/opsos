package terraform

import (
	"fmt"
	"path"
	"strings"

	"github.com/neermitt/opsos/pkg/utils"
)

func GenerateVarFileFile(ectx ExecutionContext, format string) error {
	// Write varfile to file
	var varfilePath = path.Join(ectx.WorkingDir, constructVarfileName(ectx, format))

	fmt.Printf("Writing the backend config to file:\n%s\n", varfilePath)
	if ectx.DryRun {
		return nil
	}
	return utils.PrintOrWriteToFile(format, varfilePath, ectx.ComponentConfig.Vars, 0644)
}

func constructVarfileName(ectx ExecutionContext, _ string) string {
	fmtdComponentFolderPrefix := strings.ReplaceAll(ectx.ComponentName, "/", "-")
	return fmt.Sprintf("%s-%s.terraform.tfvars.json", ectx.Stack.Name, fmtdComponentFolderPrefix)
}
