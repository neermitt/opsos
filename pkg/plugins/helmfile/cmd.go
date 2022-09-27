package helmfile

import (
	"bytes"
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

type helmfileApplyOptions struct {
	DryRun bool
}

var (
	applyOptions helmfileApplyOptions
)

// terraformCmd represents the base command for all terraform sub-commands
var helmfileApplyCmd = &cobra.Command{
	Use:                "apply stack component",
	Short:              "Execute 'helmfile apply' command",
	Long:               `This command runs helmfile apply <stack> <component>`,
	Args:               cobra.MinimumNArgs(2),
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		stackName := args[0]
		component := args[1]

		additionalArgs := args[2:]

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

		varFile := constructHelmfileComponentVarfileName(stk, info)
		varFilePath := constructHelmfileComponentVarfilePath(stk, info)

		fmt.Println("Writing the variables to file:")
		fmt.Println(varFilePath)

		if !applyOptions.DryRun {
			err = utils.PrintOrWriteToFile("yaml", varFilePath, info.Vars, 0644)
			if err != nil {
				return err
			}
		}

		// Prepare arguments and flags
		commandArgs := []string{"--state-values-file", varFile, "apply"}

		commandArgs = append(commandArgs, additionalArgs...)

		cmdEnv := append([]string{},
			fmt.Sprintf("STACK=%s", stk.Name),
		)
		conf := config.GetConfig(ctx)
		for k, v := range conf.Components.Helmfile.Envs {
			pv, err := processTemplate(v, stk.Vars)
			if err != nil {
				return err
			}
			cmdEnv = append(cmdEnv, fmt.Sprintf("%s=%s", k, pv))
		}
		return utils.ExecuteShellCommand(ctx, "helmfile", commandArgs, info.WorkingDir, cmdEnv, applyOptions.DryRun)
	},
}

type helmfileCmdProvider struct {
}

func (h *helmfileCmdProvider) RegisterCommands(cmd *cobra.Command) {
	helmfileApplyCmd.Flags().BoolVar(&applyOptions.DryRun, "dry-run", false, "run in dry run mode")
	helmfileCmd.AddCommand(helmfileApplyCmd)
	cmd.AddCommand(helmfileCmd)
}

func init() {
	plugins.RegisterCmdProvider(providerName, &helmfileCmdProvider{})
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
