package terraform

import (
	"fmt"
	"path"

	"github.com/neermitt/opsos/pkg/utils"
)

func GenerateVarFileFile(ectx ExecutionContext, format string) error {
	// Write varfile to file
	var varfilePath = path.Join(ectx.workingDir, ectx.VarFile)

	fmt.Printf("Writing the vars to file:\n%s\n", varfilePath)
	if ectx.dryRun {
		return nil
	}
	return utils.PrintOrWriteToFile(format, varfilePath, ectx.componentConfig.Vars, 0644)
}
