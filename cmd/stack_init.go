package cmd

import (
	"github.com/neermitt/opsos/internal/exec"
	"github.com/spf13/cobra"
)

// stackInit inits stack components
var stackInit = &cobra.Command{
	Use:   "init <stack>",
	Short: "Execute 'stack init' command",
	Long:  `This command inits all stack components from their configuration: opsos stack init <stack>`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]
		return exec.ExecuteStackComponentsInit(cmd.Context(), stackName, componentInitOptions)
	},
}

func init() {
	stackInit.Flags().BoolVar(&componentInitOptions.DryRun, "dry-run", false, "run in dry run mode")
	stackCmd.AddCommand(stackInit)
}
