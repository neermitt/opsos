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

type ExecutionContext struct {
	Context  context.Context
	PlanFile string
	VarFile  string

	stackName       string
	componentName   string
	workspaceName   string
	componentConfig stack.ConfigWithMetadata

	execOptions utils.ExecOptions
}

func NewExecutionContext(ctx context.Context, stackName string, component string, dryRun bool) (ExecutionContext, error) {
	stk, err := stack.LoadStack(ctx, stack.LoadStackOptions{Stack: stackName, ComponentType: ComponentType, ComponentName: component})
	if err != nil {
		return ExecutionContext{}, err
	}

	terraformComponents, found := stk.Components[ComponentType]
	if !found {
		return ExecutionContext{}, fmt.Errorf("no terraform component found")
	}
	componentConfig, found := terraformComponents[component]
	if !found {
		return ExecutionContext{}, fmt.Errorf("terraform component %s not found", component)
	}

	if err != nil {
		return ExecutionContext{}, err
	}

	conf := config.GetConfig(ctx)
	workingDir := components.GetWorkingDirectory(conf, ComponentType, componentConfig.Component)

	cmdEnv, err := buildCommandEnvs(componentConfig)
	if err != nil {
		return ExecutionContext{}, err
	}

	workspaceName, err := ConstructWorkspaceName(stk, component, componentConfig)
	if err != nil {
		return ExecutionContext{}, err
	}

	planFile := constructPlanfileName(stk, component)
	varFile := constructVarfileName(stk, component)

	exeCtx := ExecutionContext{
		Context:         ctx,
		stackName:       stk.Id,
		componentName:   component,
		componentConfig: componentConfig,
		workspaceName:   workspaceName,
		PlanFile:        planFile,
		VarFile:         varFile,
		execOptions: utils.ExecOptions{
			Env:              cmdEnv,
			WorkingDirectory: workingDir,
			StdOut:           os.Stdout,
			DryRun:           dryRun,
		},
	}
	return exeCtx, nil
}

func getCommand(exeCtx ExecutionContext) string {
	command := "terraform"
	if exeCtx.componentConfig.Command != nil {
		command = *exeCtx.componentConfig.Command
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

func ExecuteCommand(exeCtx ExecutionContext, args []string) error {
	command := getCommand(exeCtx)

	return utils.ExecuteShellCommand(exeCtx.Context, command, args, exeCtx.execOptions)
}

func ExecuteShell(exeCtx ExecutionContext) error {
	execOptions := exeCtx.execOptions
	execOptions.Env = append(execOptions.Env,
		fmt.Sprintf("TF_CLI_ARGS_plan=-var-file=%s", exeCtx.VarFile),
		fmt.Sprintf("TF_CLI_ARGS_apply=-var-file=%s", exeCtx.VarFile),
		fmt.Sprintf("TF_CLI_ARGS_refresh=-var-file=%s", exeCtx.VarFile),
		fmt.Sprintf("TF_CLI_ARGS_import=-var-file=%s", exeCtx.VarFile),
		fmt.Sprintf("TF_CLI_ARGS_destroy=-var-file=%s", exeCtx.VarFile),
		fmt.Sprintf("TF_CLI_ARGS_console=-var-file=%s", exeCtx.VarFile),
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

	return utils.ExecuteShellCommand(exeCtx.Context, shellCommand, args, execOptions)
}
