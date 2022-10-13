package exec

import (
	"context"
	"fmt"
	"github.com/neermitt/opsos/pkg/config"
	"os"
	"path"
	"strings"

	"github.com/neermitt/opsos/pkg/stack"
	"github.com/neermitt/opsos/pkg/utils"
)

type TerraformGenerateVarfileOptions struct {
	DryRun bool
	Format string
}

// ExecuteTerraformGenerateVarfile executes `terraform generate varfile` command
func ExecuteTerraformGenerateVarfile(ctx context.Context, stackName string, component string, options TerraformGenerateVarfileOptions) error {
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

	fmt.Print("Component backend config:\\n\\n")
	err = utils.Get("json")(os.Stdout, info.Vars)
	if err != nil {
		return err
	}

	conf := config.GetConfig(ctx)
	workingDir, _, err := getComponentWorkingDirectory(conf, terraformComponentType, info)
	if err != nil {
		return err
	}

	// Write backend config to file
	var backendFilePath = constructTerraformComponentVarfilePath(stk, component, workingDir)

	fmt.Println()
	fmt.Printf("Writing the backend config to file:\n%s\n", backendFilePath)
	if !options.DryRun {
		err = utils.PrintOrWriteToFile(options.Format, backendFilePath, info.Vars, 0644)
		if err != nil {
			return err
		}

	}
	return nil
}

func constructTerraformComponentVarfileName(stk *stack.Stack, componentName string) string {
	fmtdComponentFolderPrefix := strings.ReplaceAll(componentName, "/", "-")
	return fmt.Sprintf("%s-%s.terraform.tfvars.json", stk.Name, fmtdComponentFolderPrefix)
}

func constructTerraformComponentVarfilePath(stk *stack.Stack, componentName string, workingDir string) string {
	return path.Join(workingDir, constructTerraformComponentVarfileName(stk, componentName))
}
