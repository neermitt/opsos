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
	Context  context.Context
	Config   *config.Configuration
	PlanFile string
	VarFile  string

	dryRun bool

	stackName       string
	componentName   string
	componentConfig stack.ConfigWithMetadata
	workingDir      string
	workspaceName   string
	processedCmdEnv []string
}

type Option func(execCtx *ExecutionContext)

func WithDryRun(enable bool) Option {
	return func(execCtx *ExecutionContext) {
		execCtx.dryRun = enable
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
		stackName:       stk.Id,
		componentName:   component,
		componentConfig: componentConfig,
		workingDir:      workingDir,
		workspaceName:   workspaceName,
		PlanFile:        planFile,
		VarFile:         varFile,
		processedCmdEnv: cmdEnv,
	}
	for _, opt := range options {
		opt(&exeCtx)
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

	return utils.ExecuteShellCommand(exeCtx.Context, command, args, utils.ExecOptions{
		DryRun:           exeCtx.dryRun,
		Env:              exeCtx.processedCmdEnv,
		WorkingDirectory: exeCtx.workingDir,
		StdOut:           os.Stdout,
	})
}
