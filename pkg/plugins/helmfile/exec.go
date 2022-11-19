package helmfile

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/neermitt/opsos/pkg/components"
	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/plugins"
	"github.com/neermitt/opsos/pkg/stack"
	"github.com/neermitt/opsos/pkg/utils"
)

func ExecHelmfileCommand(ctx context.Context, command string, stk *stack.Stack, globalArgs []string, additionalArgs []string, dryRun bool) error {
	stackName := stack.GetStackName(ctx)
	component := stack.GetComponent(ctx)
	helmfileCM, found := stk.Components[component.Type]
	if !found {
		return fmt.Errorf("not helmfile component found")
	}
	info, found := helmfileCM[component.Name]
	if !found {
		return fmt.Errorf("helmfile component %s not found", component)
	}

	// Export KubeConfig
	conf := config.GetConfig(ctx)

	var helmfileConfig Config
	err := utils.FromMap(conf.Providers[ComponentType], &helmfileConfig)
	if err != nil {
		return err
	}

	kubeConfigPath := fmt.Sprintf("%s/%s-kubecfg", helmfileConfig.KubeconfigPath, stk.Name)
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
	if err := plugins.GetKubeConfig(ctx, k8sProvider, stk, kubeConfigPath); err != nil {
		return err
	}

	// Working Dir
	workingDir := components.GetWorkingDirectory(conf, component.Type, component.Name)

	varFile := constructHelmfileComponentVarfileName(stk, component.Name)
	varFilePath := constructHelmfileComponentVarfilePath(stk, workingDir, component.Name)

	fmt.Println("Writing the variables to file:")
	fmt.Println(varFilePath)

	if !dryRun {
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
		fmt.Sprintf("KUBECONFIG=%s", kubeConfigPath),
	)
	for k, v := range helmfileConfig.Envs {
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
		DryRun:           dryRun,
		Env:              cmdEnv,
		WorkingDirectory: workingDir,
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
