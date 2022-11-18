package terraform

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/neermitt/opsos/pkg/components"
	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/stack"
	"github.com/neermitt/opsos/pkg/utils"
)

type Settings struct {
	PlanFile      string
	VarFile       string
	WorkspaceName string
}

func setTerraformSettings(ctx context.Context, config Settings) context.Context {
	return context.WithValue(ctx, "terraform-settings", config)
}

func GetTerraformSettings(ctx context.Context) Settings {
	return ctx.Value("terraform-settings").(Settings)
}

func NewExecutionContext(ctx context.Context, stk *stack.Stack, component stack.Component, dryRun bool) (context.Context, error) {
	terraformComponents, found := stk.Components[component.Type]
	if !found {
		return nil, fmt.Errorf("no terraform component found")
	}
	componentConfig, found := terraformComponents[component.Name]
	if !found {
		return nil, fmt.Errorf("terraform component %s not found", component)
	}

	conf := config.GetConfig(ctx)
	workingDir := components.GetWorkingDirectory(conf, ComponentType, componentConfig.Component)

	cmdEnv, err := buildCommandEnvs(componentConfig)
	if err != nil {
		return nil, err
	}

	workspaceName, err := ConstructWorkspaceName(stk, component.Name, componentConfig)
	if err != nil {
		return nil, err
	}

	planFile := constructPlanfileName(stk, component.Name)
	varFile := constructVarfileName(stk, component.Name)

	ctx = stack.SetComponentConfig(ctx, &componentConfig)

	ctx = setTerraformSettings(ctx, Settings{
		PlanFile:      planFile,
		VarFile:       varFile,
		WorkspaceName: workspaceName,
	})

	ctx = utils.SetExecOptions(ctx, utils.ExecOptions{
		Env:              cmdEnv,
		WorkingDirectory: workingDir,
		StdOut:           os.Stdout,
		DryRun:           dryRun,
	})

	return ctx, nil
}

func getCommand(ctx context.Context) string {
	command := "terraform"
	componentConfig := stack.GetComponentConfig(ctx)
	if componentConfig.Command != nil {
		command = *componentConfig.Command
	}
	return command
}

func buildCommandEnvs(config stack.ConfigWithMetadata) ([]string, error) {
	var cmdEnv []string
	for k, v := range config.Envs {
		pv, err := utils.ProcessTemplate(v, config.Vars)
		if err != nil {
			return nil, err
		}
		cmdEnv = append(cmdEnv, fmt.Sprintf("%s=%s", k, pv))
	}
	return cmdEnv, nil
}

func ExecuteCommand(ctx context.Context, args []string) error {
	command := getCommand(ctx)
	execOptions := utils.GetExecOptions(ctx)

	return utils.ExecuteShellCommand(ctx, command, args, execOptions)
}

func ExecuteShell(ctx context.Context) error {
	terraformOptions := GetTerraformSettings(ctx)
	execOptions := utils.GetExecOptions(ctx)
	execOptions.Env = append(execOptions.Env,
		fmt.Sprintf("TF_CLI_ARGS_plan=-var-file=%s", terraformOptions.VarFile),
		fmt.Sprintf("TF_CLI_ARGS_apply=-var-file=%s", terraformOptions.VarFile),
		fmt.Sprintf("TF_CLI_ARGS_refresh=-var-file=%s", terraformOptions.VarFile),
		fmt.Sprintf("TF_CLI_ARGS_import=-var-file=%s", terraformOptions.VarFile),
		fmt.Sprintf("TF_CLI_ARGS_destroy=-var-file=%s", terraformOptions.VarFile),
		fmt.Sprintf("TF_CLI_ARGS_console=-var-file=%s", terraformOptions.VarFile),
	)

	var shellCommand string
	var args []string
	if runtime.GOOS == "windows" {
		shellCommand = "cmd.exe"
	} else {
		// If 'SHELL' ENV var is not defined, use 'bash' shell
		shellCommand = os.Getenv("SHELL")
		if len(shellCommand) == 0 {
			bashPath, err := exec.LookPath("bash")
			if err != nil {
				return err
			}
			shellCommand = bashPath
		}
		args = append(args, "-l")
	}

	return utils.ExecuteShellCommand(ctx, shellCommand, args, execOptions)
}
