package cmd

import (
	"github.com/neermitt/opsos/internal/exec"
	"github.com/spf13/cobra"
)

var (
	componentInitOptions exec.ComponentInitOptions
)

// componentInitCmd inits components
var componentInitCmd = &cobra.Command{
	Use:   "init <component-type> <component>",
	Short: "Execute 'component init' command",
	Long:  `This command inits component from its configuration: opsos component init <component-type> <component>`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		componentInitOptions.ComponentType = args[0]
		componentInitOptions.ComponentName = args[1]
		return exec.ExecuteComponentInit(cmd.Context(), componentInitOptions)
	},
}

func init() {
	componentInitCmd.Flags().BoolVar(&componentInitOptions.DryRun, "dry-run", false, "run in dry run mode")
	componentCmd.AddCommand(componentInitCmd)
}
