package exec

import (
	"context"
	"errors"
	"os"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/plugins/terraform"
	"github.com/neermitt/opsos/pkg/stack"
	"github.com/neermitt/opsos/pkg/utils"
)

type TerraformOptions struct {
	Command                   string
	RequiresVarFile           bool
	UsePlan                   bool
	DryRun                    bool
	AutoApprove               bool
	Destroy                   bool
	SkipInit                  bool
	SkipWorkspace             bool
	CleanPlanFileOnCompletion bool
}

func ExecuteTerraform(ctx context.Context, stackName string, componentName string, additionalArgs []string, options TerraformOptions) error {

	component := stack.Component{Type: terraform.ComponentType, Name: componentName}
	ctx = stack.SetStackName(ctx, stackName)
	ctx = stack.SetComponent(ctx, component)
	stk, err := stack.LoadStack(ctx, stack.LoadStackOptions{Stack: stackName, Component: &component})
	if err != nil {
		return err
	}
	ctx, err = terraform.NewExecutionContext(ctx, stk, component, options.DryRun)
	if err != nil {
		return err
	}

	if options.RequiresVarFile {
		if err := terraform.GenerateVarFileFile(ctx, "json"); err != nil {
			return err
		}
	}

	conf := config.GetConfig(ctx)
	var terraformConfig terraform.Config
	err = utils.FromMap(conf.Providers[terraform.ComponentType], &terraformConfig)
	if err != nil {
		return err
	}

	if terraformConfig.AutoGenerateBackendFile {
		if err := terraform.GenerateBackendFile(ctx, "json"); err != nil {
			return err
		}
	}

	if !options.SkipInit {
		initArgs, err := prepareArgs(ctx, terraformConfig, nil, TerraformOptions{Command: "init"})
		if err != nil {
			return err
		}
		if err := terraform.ExecuteCommand(ctx, initArgs); err != nil {
			return err
		}
	}
	if !options.SkipWorkspace {
		if err := terraform.SelectOrCreateWorkspace(ctx); err != nil {
			return err
		}
	}

	args, err := prepareArgs(ctx, terraformConfig, additionalArgs, options)
	if err != nil {
		return err
	}

	if options.Command == "shell" {
		return terraform.ExecuteShell(ctx)
	}
	if err := terraform.ExecuteCommand(ctx, args); err != nil {
		return err
	}
	if options.CleanPlanFileOnCompletion {
		return terraform.CleanPlanFile(ctx)
	}
	return nil
}

func prepareArgs(ctx context.Context, terraformConfig terraform.Config, additionalArgs []string, options TerraformOptions) ([]string, error) {
	args := []string{options.Command}

	terraformSettings := terraform.GetTerraformSettings(ctx)

	if options.UsePlan {
		args = append(args, terraformSettings.PlanFile)
	} else if options.RequiresVarFile {
		args = append(args, "-var-file", terraformSettings.VarFile)
	}

	switch options.Command {
	case "plan":
		args = append(args, "-out", terraformSettings.PlanFile)
	case "init":
		if terraformConfig.InitRunReconfigure {
			args = append(args, "-reconfigure")
		}
	case "apply":
		if !utils.StringInSlice(terraform.AutoApproveFlag, additionalArgs) {
			if (terraformConfig.ApplyAutoApprove || options.AutoApprove) && !options.UsePlan {
				args = append(args, terraform.AutoApproveFlag)
			} else if os.Stdin == nil {
				return nil, errors.New("'terraform apply' requires a user interaction, but it's running without `tty` or `stdin` attached.\nUse 'terraform apply -auto-approve' or 'terraform deploy' instead.")
			}
		}
		if options.Destroy {
			args = append(args, "-destroy")
		}
	}
	args = append(args, additionalArgs...)
	return args, nil
}
