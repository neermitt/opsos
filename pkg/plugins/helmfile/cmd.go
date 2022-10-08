package helmfile

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/plugins"
	"github.com/neermitt/opsos/pkg/stack"
	"github.com/neermitt/opsos/pkg/utils"
	"github.com/spf13/cobra"
)

// terraformCmd represents the base command for all terraform sub-commands
var helmfileCmd = &cobra.Command{
	Use:   "helmfile",
	Short: "Execute 'helmfile' commands",
	Long:  `This command runs helmfile commands`,
}

var helmfileApplyCmd = &cobra.Command{
	Use:                "apply stack component",
	Short:              "Execute 'helmfile apply' command",
	Long:               `This command runs helmfile apply <stack> <component>`,
	Args:               cobra.MinimumNArgs(2),
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]
		component := args[1]

		additionalArgs := args[2:]

		return execHelmfile(cmd.Context(), "apply", stackName, component, additionalArgs)
	},
}

var helmfileSyncCmd = &cobra.Command{
	Use:                "sync stack component",
	Short:              "Execute 'helmfile sync' command",
	Long:               `This command runs helmfile sync <stack> <component>`,
	Args:               cobra.MinimumNArgs(2),
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]
		component := args[1]

		additionalArgs := args[2:]

		return execHelmfile(cmd.Context(), "sync", stackName, component, additionalArgs)
	},
}

var helmfileDiffCmd = &cobra.Command{
	Use:                "diff stack component",
	Short:              "Execute 'helmfile diff' command",
	Long:               `This command runs helmfile diff <stack> <component>`,
	Args:               cobra.MinimumNArgs(2),
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]
		component := args[1]

		additionalArgs := args[2:]

		return execHelmfile(cmd.Context(), "diff", stackName, component, additionalArgs)
	},
}

type helmfileCmdProvider struct {
}

func (h *helmfileCmdProvider) RegisterCommands(cmd *cobra.Command) {
	helmfileApplyCmd.Flags().BoolVar(&helmfileExecOptions.DryRun, "dry-run", false, "run in dry run mode")
	helmfileApplyCmd.Flags().StringVar(&helmfileExecOptions.GlobalArgs, "global-args", "", "global options of `helmfile`")

	helmfileSyncCmd.Flags().BoolVar(&helmfileExecOptions.DryRun, "dry-run", false, "run in dry run mode")
	helmfileSyncCmd.Flags().StringVar(&helmfileExecOptions.GlobalArgs, "global-args", "", "global options of `helmfile`")

	helmfileDiffCmd.Flags().StringVar(&helmfileExecOptions.GlobalArgs, "global-args", "", "global options of `helmfile`")

	helmfileCmd.AddCommand(helmfileApplyCmd)
	helmfileCmd.AddCommand(helmfileSyncCmd)
	helmfileCmd.AddCommand(helmfileDiffCmd)

	cmd.AddCommand(helmfileCmd)
}

func init() {
	plugins.RegisterCmdProvider(providerName, &helmfileCmdProvider{})
}

var (
	helmfileExecOptions struct {
		DryRun     bool
		GlobalArgs string
	}
)

func execHelmfile(ctx context.Context, command string, stackName string, component string, additionalArgs []string) error {

	globalArgs := strings.Fields(helmfileExecOptions.GlobalArgs)

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

	info := hc.(helmfileComponentInfo)

	conf := config.GetConfig(ctx)
	kubeconfigPath := fmt.Sprintf("%s/%s-kubecfg", conf.Components.Helmfile.KubeconfigPath, stk.Name)
	clusterName, err := processTemplate(conf.Components.Helmfile.ClusterNamePattern, stk.Vars)
	if err != nil {
		return err
	}
	kubeConfigProvider, found := plugins.GetKubeConfigProvider(ctx, stk.KubeConfigProvider)
	if !found {
		return fmt.Errorf("%s kube config provider is not configured", stk.KubeConfigProvider)
	}
	err = kubeConfigProvider.ExportKubeConfig(ctx, clusterName, kubeconfigPath)
	if err != nil {
		return err
	}

	varFile := constructHelmfileComponentVarfileName(stk, info)
	varFilePath := constructHelmfileComponentVarfilePath(stk, info)

	fmt.Println("Writing the variables to file:")
	fmt.Println(varFilePath)

	if !helmfileExecOptions.DryRun {
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
		pv, err := processTemplate(v, stk.Vars)
		if err != nil {
			return err
		}
		cmdEnv = append(cmdEnv, fmt.Sprintf("%s=%s", k, pv))
	}
	return utils.ExecuteShellCommand(ctx, "helmfile", commandArgs, utils.ExecOptions{
		DryRun:           helmfileExecOptions.DryRun,
		Env:              cmdEnv,
		WorkingDirectory: info.WorkingDir,
	})
}

func constructHelmfileComponentVarfileName(stk *stack.Stack, info helmfileComponentInfo) string {
	var varFile string
	if len(info.ComponentFolderPrefix) == 0 {
		varFile = fmt.Sprintf("%s-%s.helmfile.vars.yaml", stk.Name, info.Component)
	} else {
		varFile = fmt.Sprintf("%s-%s-%s.helmfile.vars.yaml", stk.Name, strings.ReplaceAll(info.ComponentFolderPrefix, "/", "-"), info.Component)
	}
	return varFile
}

// constructHelmfileComponentVarfilePath constructs the varfile path for a helmfile component in a stack
func constructHelmfileComponentVarfilePath(stk *stack.Stack, info helmfileComponentInfo) string {
	return path.Join(
		info.WorkingDir,
		constructHelmfileComponentVarfileName(stk, info),
	)
}

func processTemplate(s string, vars map[string]any) (string, error) {
	tmpl, err := template.New("template").Parse(s)
	if err != nil {
		return "", err
	}
	var buff bytes.Buffer
	if err := tmpl.Execute(&buff, vars); err != nil {
		return "", err
	}
	return buff.String(), nil
}
