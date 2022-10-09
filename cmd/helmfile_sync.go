package cmd

import (
	"github.com/neermitt/opsos/internal/exec"
	"github.com/spf13/cobra"
)

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

		return exec.ExecHelmfile(cmd.Context(), "sync", stackName, component, additionalArgs, helmfileExecOptions)
	},
}

func init() {
	helmfileSyncCmd.Flags().BoolVar(&helmfileExecOptions.DryRun, "dry-run", false, "run in dry run mode")
	helmfileSyncCmd.Flags().StringVar(&helmfileExecOptions.GlobalArgs, "global-args", "", "global options of `helmfile`")
	helmfileCmd.AddCommand(helmfileSyncCmd)
}
