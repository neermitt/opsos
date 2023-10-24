package cmds

import (
	"github.com/neermitt/opsos/pkg/plugins/helmfile/exec"
	"github.com/spf13/cobra"
)

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

		return exec.ExecHelmfile(cmd.Context(), "apply", stackName, component, additionalArgs, helmfileExecOptions)
	},
}

func init() {
	helmfileApplyCmd.Flags().BoolVar(&helmfileExecOptions.DryRun, "dry-run", false, "run in dry run mode")
	helmfileApplyCmd.Flags().StringVar(&helmfileExecOptions.GlobalArgs, "global-args", "", "global options of `helmfile`")
	helmfileCmd.AddCommand(helmfileApplyCmd)
}
