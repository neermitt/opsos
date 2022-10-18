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

	stackName       string
	componentName   string
	componentConfig stack.ConfigWithMetadata
	workspaceName   string

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
		Config:          conf,
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
