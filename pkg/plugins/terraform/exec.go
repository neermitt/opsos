package terraform

import (
	"context"
	"fmt"

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
	DryRun          bool
	AdditionalArgs  []string
	PlanFile        string
}

func getCommand(exeCtx ExecutionContext) string {
	command := "terraform"
	if exeCtx.ComponentConfig.Command != nil {
		command = *exeCtx.ComponentConfig.Command
	}
	return command
}

func buildCommandEnvs(exeCtx ExecutionContext) ([]string, error) {
	var cmdEnv []string
	for k, v := range exeCtx.ComponentConfig.Envs {
		pv, err := utils.ProcessTemplate(v, exeCtx.ComponentConfig.Vars)
		if err != nil {
			return nil, err
		}
		cmdEnv = append(cmdEnv, fmt.Sprintf("%s=%s", k, pv))
	}
	return cmdEnv, nil
}
