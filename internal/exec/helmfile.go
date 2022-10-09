package exec

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/plugins"
	"github.com/neermitt/opsos/pkg/plugins/helmfile"
	"github.com/neermitt/opsos/pkg/stack"
	"github.com/neermitt/opsos/pkg/utils"
)

const (
	providerName = "helmfile"
)

type HelmfileExecOptions struct {
	DryRun     bool
	GlobalArgs string
}

func ExecHelmfile(ctx context.Context, command string, stackName string, component string, additionalArgs []string, options HelmfileExecOptions) error {

	globalArgs := strings.Fields(options.GlobalArgs)

	stk, err := stack.LoadStack(ctx, stack.LoadStackOptions{Stack: stackName})
	if err != nil {
		return err
	}

	helmfileCM, found := stk.ComponentTypes[providerName]
	if !found {
		return fmt.Errorf("not helmfile component found")
	}
	hc, found := helmfileCM[component]
	if !found {
		return fmt.Errorf("helmfile component %s not found", component)
	}

	info := hc.(helmfile.ComponentInfo)

	conf := config.GetConfig(ctx)
	kubeconfigPath := fmt.Sprintf("%s/%s-kubecfg", conf.Components.Helmfile.KubeconfigPath, stk.Name)
	clusterName, err := utils.ProcessTemplate(conf.Components.Helmfile.ClusterNamePattern, stk.Vars)
	if err != nil {
		return err
	}
	if err := plugins.GetKubeConfig(ctx, stk.KubeConfigProvider, clusterName, kubeconfigPath); err != nil {
		return err
	}

	varFile := constructHelmfileComponentVarfileName(stk, info)
	varFilePath := constructHelmfileComponentVarfilePath(stk, info)

	fmt.Println("Writing the variables to file:")
	fmt.Println(varFilePath)

	if !options.DryRun {
		err = utils.PrintOrWriteToFile("yaml", varFilePath, info.Vars, 0644)
		if err != nil {
			return err
		}
	}

	// Prepare arguments and flags
	commandArgs := []string{"--state-values-file", varFile}

	commandArgs = append(commandArgs, globalArgs...)

	commandArgs = append(commandArgs, command)

	commandArgs = append(commandArgs, additionalArgs...)

	cmdEnv := append([]string{},
		fmt.Sprintf("STACK=%s", stk.Name),
		fmt.Sprintf("KUBECONFIG=%s", kubeconfigPath),
	)
	for k, v := range conf.Components.Helmfile.Envs {
		pv, err := utils.ProcessTemplate(v, stk.Vars)
		if err != nil {
			return err
		}
		cmdEnv = append(cmdEnv, fmt.Sprintf("%s=%s", k, pv))
	}
	return utils.ExecuteShellCommand(ctx, "helmfile", commandArgs, utils.ExecOptions{
		DryRun:           options.DryRun,
		Env:              cmdEnv,
		WorkingDirectory: info.WorkingDir,
	})
}

func constructHelmfileComponentVarfileName(stk *stack.Stack, info helmfile.ComponentInfo) string {
	var varFile string
	if len(info.ComponentFolderPrefix) == 0 {
		varFile = fmt.Sprintf("%s-%s.helmfile.vars.yaml", stk.Name, info.Component)
	} else {
		fmtdComponentFolderPrefix := strings.ReplaceAll(info.ComponentFolderPrefix, "/", "-")
		varFile = fmt.Sprintf("%s-%s-%s.helmfile.vars.yaml", stk.Name, fmtdComponentFolderPrefix, info.Component)
	}
	return varFile
}

// constructHelmfileComponentVarfilePath constructs the varfile path for a helmfile component in a stack
func constructHelmfileComponentVarfilePath(stk *stack.Stack, info helmfile.ComponentInfo) string {
	return path.Join(
		info.WorkingDir,
		constructHelmfileComponentVarfileName(stk, info),
	)
}
