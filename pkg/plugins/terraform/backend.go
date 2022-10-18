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

	if ectx.ComponentConfig.BackendType == nil || *ectx.ComponentConfig.BackendType == "" {
		return fmt.Errorf("'backend_type' is missing for the '%[2]s' component in stack %[1]s", ectx.Stack.Id, ectx.ComponentName)
	}

	if ectx.ComponentConfig.Backend == nil {
		return fmt.Errorf("could not find 'backend' config for the '%[2]s' component in stack %[1]s", ectx.Stack.Id, ectx.ComponentName)
	}

	fmt.Print("Component backend config:\n\n")
	err := utils.GetFormatter("json")(os.Stdout, ectx.ComponentConfig.Backend)
	if err != nil {
		return err
	}

	componentBackendConfig := Root{
		Terraform{
			Backend: Backend{
				Type: *ectx.ComponentConfig.BackendType,
				Data: ectx.ComponentConfig.Backend,
			},
			JSONBackend: map[string]map[string]any{*ectx.ComponentConfig.BackendType: ectx.ComponentConfig.Backend},
		}}

	// Write backend config to file
	var backendFilePath = path.Join(ectx.WorkingDir, constructBackendFileName(ectx, format))

	fmt.Printf("Writing the backend config to file:\n%s\n", backendFilePath)
	if ectx.DryRun {
		return nil
	}
	return utils.PrintOrWriteToFile(format, backendFilePath, componentBackendConfig, 0644)
}

// constructBackendFileName constructs the backend path for a terraform component in a stack
func constructBackendFileName(_ ExecutionContext, format string) string {
	if format == "json" {
		return "backend.tf.json"
	}
	return "backend.tf"
}
