package terraform

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/neermitt/opsos/pkg/stack"
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

func GenerateBackendFile(ctx context.Context, format string) error {
	componentConfig := stack.GetComponentConfig(ctx)

	if componentConfig.BackendType == nil || *componentConfig.BackendType == "" {
		stackName := stack.GetStackName(ctx)
		component := stack.GetComponent(ctx)
		return fmt.Errorf("'backend_type' is missing for the '%[2]s' component in stack %[1]s", stackName, component.Name)
	}

	if componentConfig.Backend == nil {
		stackName := stack.GetStackName(ctx)
		component := stack.GetComponent(ctx)
		return fmt.Errorf("could not find 'backend' config for the '%[2]s' component in stack %[1]s", stackName, component.Name)
	}

	fmt.Print("Component backend config:\n\n")
	err := utils.GetFormatter("json")(os.Stdout, componentConfig.Backend)
	if err != nil {
		return err
	}

	componentBackendConfig := Root{
		Terraform{
			Backend: Backend{
				Type: *componentConfig.BackendType,
				Data: componentConfig.Backend,
			},
			JSONBackend: map[string]map[string]any{*componentConfig.BackendType: componentConfig.Backend},
		}}

	execOptions := utils.GetExecOptions(ctx)

	// Write backend config to file
	var backendFilePath = path.Join(execOptions.WorkingDirectory, constructBackendFileName(format))

	log.Printf("[INFO] (terraform) Writing the backend config to file: %s", backendFilePath)
	if execOptions.DryRun {
		return nil
	}
	return utils.PrintOrWriteToFile(format, backendFilePath, componentBackendConfig, 0644)
}
