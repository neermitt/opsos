package exec

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/plugins"
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

	stk, err := stack.LoadStack(ctx, stack.LoadStackOptions{Stack: stackName, ComponentType: providerName, ComponentName: component})
	if err != nil {
		return err
	}

	helmfileCM, found := stk.Components[providerName]
	if !found {
		return fmt.Errorf("not helmfile component found")
	}
	info, found := helmfileCM[component]
	if !found {
		return fmt.Errorf("helmfile component %s not found", component)
	}

	// Export KubeConfig
	conf := config.GetConfig(ctx)
	kubeconfigPath := fmt.Sprintf("%s/%s-kubecfg", conf.Components.Helmfile.KubeconfigPath, stk.Name)
	clusterName, err := utils.ProcessTemplate(conf.Components.Helmfile.ClusterNamePattern, stk.Vars)
	if err != nil {
		return err
	}

	k8sSettings, found := info.Settings["k8s"].(map[string]any)
	if !found {
		return fmt.Errorf("k8s settings not found for component %[2]s in stack %[1]s", stackName, component)
	}
	k8sProvider, found := k8sSettings["provider"].(string)
	if !found {
		return fmt.Errorf("k8s settings for provider not found for component %[2]s in stack %[1]s", stackName, component)
	}
	if err := plugins.GetKubeConfig(ctx, k8sProvider, clusterName, kubeconfigPath); err != nil {
		return err
	}

	// Working Dir
	workingDir := path.Join(conf.BasePath, conf.Components.Helmfile.BasePath, info.Component)
	absWorkingDir, err := filepath.Abs(workingDir)
	if err != nil {
		return err
	}

	varFile := constructHelmfileComponentVarfileName(stk, component)
	varFilePath := constructHelmfileComponentVarfilePath(stk, workingDir, component)

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
	for k, v := range info.Envs {
		pv, err := utils.ProcessTemplate(v, stk.Vars)
		if err != nil {
			return err
		}
		cmdEnv = append(cmdEnv, fmt.Sprintf("%s=%s", k, pv))
	}

	return utils.ExecuteShellCommand(ctx, "helmfile", commandArgs, utils.ExecOptions{
		DryRun:           options.DryRun,
		Env:              cmdEnv,
		WorkingDirectory: absWorkingDir,
	})
}

func constructHelmfileComponentVarfileName(stk *stack.Stack, componentName string) string {
	fmtdComponentFolderPrefix := strings.ReplaceAll(componentName, "/", "-")
	return fmt.Sprintf("%s-%s.helmfile.vars.yaml", stk.Name, fmtdComponentFolderPrefix)
}

// constructHelmfileComponentVarfilePath constructs the varfile path for a helmfile component in a stack
func constructHelmfileComponentVarfilePath(stk *stack.Stack, workingDir string, componentName string) string {
	return path.Join(
		workingDir,
		constructHelmfileComponentVarfileName(stk, componentName),
	)
}
