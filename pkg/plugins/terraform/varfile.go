package terraform

import (
	"fmt"
	"github.com/neermitt/opsos/pkg/utils"
	"path"
)

func GenerateVarFileFile(ectx ExecutionContext, format string) error {
	// Write varfile to file
	var varfilePath = path.Join(ectx.WorkingDir, constructVarfileName(ectx.Stack, ectx.ComponentName))

	fmt.Printf("Writing the vars to file:\n%s\n", varfilePath)
	if ectx.DryRun {
		return nil
	}
	return utils.PrintOrWriteToFile(format, varfilePath, ectx.ComponentConfig.Vars, 0644)
}
