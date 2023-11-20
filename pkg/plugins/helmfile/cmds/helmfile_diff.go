package cmds

import (
	"github.com/neermitt/opsos/pkg/plugins/helmfile/exec"
	"github.com/spf13/cobra"
)

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

		return exec.ExecHelmfile(cmd.Context(), "diff", stackName, component, additionalArgs, helmfileExecOptions)
	},
}

func init() {
	helmfileDiffCmd.Flags().StringVar(&helmfileExecOptions.GlobalArgs, "global-args", "", "global options of `helmfile`")
	helmfileCmd.AddCommand(helmfileDiffCmd)
}
