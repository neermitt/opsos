package exec

import (
	"context"
	"fmt"
	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/plugins/terraform"
	"os"
	"path"

	"github.com/neermitt/opsos/pkg/stack"
	"github.com/neermitt/opsos/pkg/utils"
)

const (
	terraformComponentType = "terraform"
)

type TerraformGenerateBackendOptions struct {
	DryRun bool
	Format string
}

// ExecuteTerraformGenerateBackend executes `terraform generate backend` command
func ExecuteTerraformGenerateBackend(ctx context.Context, stackName string, component string, options TerraformGenerateBackendOptions) error {
	stk, err := stack.LoadStack(ctx, stack.LoadStackOptions{Stack: stackName, ComponentType: "terraform", ComponentName: component})
	if err != nil {
		return err
	}

	terraformComponents, found := stk.Components[terraformComponentType]
	if !found {
		return fmt.Errorf("no terraform component found")
	}
	info, found := terraformComponents[component]
	if !found {
		return fmt.Errorf("terraform component %s not found", component)
	}

	if info.BackendType == nil || *info.BackendType == "" {
		return fmt.Errorf("\n'backend_type' is missing for the '%s' component.\n", component)
	}

	if info.Backend == nil {
		return fmt.Errorf("\nCould not find 'backend' config for the '%s' component.\n", component)
	}

	var r any = terraform.GetBackend(*info.BackendType, info.Backend)
	componentBackendConfig := r

	fmt.Print("Component backend config:\\n\\n")
	err = utils.Get("json")(os.Stdout, componentBackendConfig)
	if err != nil {
		return err
	}

	conf := config.GetConfig(ctx)
	workingDir, _, err := getComponentWorkingDirectory(conf, terraformComponentType, info)
	if err != nil {
		return err
	}

	// Write backend config to file
	var backendFilePath = constructTerraformBackendfilePath(workingDir, options.Format)

	fmt.Println()
	fmt.Printf("Writing the backend config to file:\n%s\n", backendFilePath)
	if !options.DryRun {
		err = utils.PrintOrWriteToFile(options.Format, backendFilePath, componentBackendConfig, 0644)
		if err != nil {
			return err
		}

	}
	return nil
}

// constructTerraformBackendfilePath constructs the backend path for a terraform component in a stack
func constructTerraformBackendfilePath(workingDir string, format string) string {
	if format == "json" {
		return path.Join(workingDir, "backend.tf.json")
	}
	return path.Join(workingDir, "backend.tf")
}
