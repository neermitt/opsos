package terraform

import (
	"context"
	"fmt"
	"os"

	"github.com/neermitt/opsos/pkg/components"
	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/stack"
	"github.com/neermitt/opsos/pkg/utils"
)

type ExecutionContext struct {
	Context         context.Context
	Config          *config.Configuration
	Stack           *stack.Stack
	ComponentName   string
	ComponentConfig stack.ConfigWithMetadata
	WorkingDir      string
	WorkspaceName   string
	DryRun          bool
	CmdEnv          []string
	PlanFile        string
	VarFile         string
}

type Option func(execCtx *ExecutionContext)

func WithDryRun() Option {
	return func(execCtx *ExecutionContext) {
		execCtx.DryRun = true
	}
}

func NewExecutionContext(ctx context.Context, stackName string, component string, options ...Option) (ExecutionContext, error) {
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
		Config:          conf,
		Stack:           stk,
		ComponentName:   component,
		ComponentConfig: componentConfig,
		WorkingDir:      workingDir,
		WorkspaceName:   workspaceName,
		DryRun:          false,
		PlanFile:        planFile,
		VarFile:         varFile,
		CmdEnv:          cmdEnv,
	}
	for _, opt := range options {
		opt(&exeCtx)
	}
	return exeCtx, nil
}

func getCommand(exeCtx ExecutionContext) string {
	command := "terraform"
	if exeCtx.ComponentConfig.Command != nil {
		command = *exeCtx.ComponentConfig.Command
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

	return utils.ExecuteShellCommand(exeCtx.Context, command, args, utils.ExecOptions{
		DryRun:           exeCtx.DryRun,
		Env:              exeCtx.CmdEnv,
		WorkingDirectory: exeCtx.WorkingDir,
		StdOut:           os.Stdout,
	})
}
