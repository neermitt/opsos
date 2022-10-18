package terraform

import (
	"fmt"
	"os"
	"path"

	"github.com/neermitt/opsos/pkg/utils"
)

type Root struct {
	Terraform Terraform `json:"terraform" hcle:"terraform,block"`
}

type Terraform struct {
	Backend Backend `json:"-" hcle:"backend,block"`

	JSONBackend map[string]map[string]any `json:"backend,inline" `
}

type Backend struct {
	Type string         `hcle:",label"`
	Data map[string]any `hcle:",body"`
}

func GenerateBackendFile(ectx ExecutionContext, format string) error {

	if ectx.componentConfig.BackendType == nil || *ectx.componentConfig.BackendType == "" {
		return fmt.Errorf("'backend_type' is missing for the '%[2]s' component in stack %[1]s", ectx.stackName, ectx.componentName)
	}

	if ectx.componentConfig.Backend == nil {
		return fmt.Errorf("could not find 'backend' config for the '%[2]s' component in stack %[1]s", ectx.stackName, ectx.componentName)
	}

	fmt.Print("Component backend config:\n\n")
	err := utils.GetFormatter("json")(os.Stdout, ectx.componentConfig.Backend)
	if err != nil {
		return err
	}

	componentBackendConfig := Root{
		Terraform{
			Backend: Backend{
				Type: *ectx.componentConfig.BackendType,
				Data: ectx.componentConfig.Backend,
			},
			JSONBackend: map[string]map[string]any{*ectx.componentConfig.BackendType: ectx.componentConfig.Backend},
		}}

	// Write backend config to file
	var backendFilePath = path.Join(ectx.workingDir, constructBackendFileName(format))

	fmt.Printf("Writing the backend config to file:\n%s\n", backendFilePath)
	if ectx.dryRun {
		return nil
	}
	return utils.PrintOrWriteToFile(format, backendFilePath, componentBackendConfig, 0644)
}
