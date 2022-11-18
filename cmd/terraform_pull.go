package cmd

import (
	"github.com/neermitt/opsos/internal/exec"
	"github.com/neermitt/opsos/pkg/plugins/terraform"
	"github.com/spf13/cobra"
)

var (
	componentPullOptions exec.ComponentPullOptions
)

// terraformPullCmd pulls components from configuration
var terraformPullCmd = &cobra.Command{
	Use:   "pull <stack> <component>",
	Short: "Execute 'terraform pull' commands",
	Long:  `This command pull component from its configuration: opsos terraform pull <stack> <component>`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]
		component := args[1]
		componentPullOptions.ComponentType = terraform.ComponentType
		return exec.ExecuteComponentPull(cmd.Context(), stackName, component, componentPullOptions)
	},
}

func init() {
	terraformPullCmd.Flags().BoolVar(&terraformOptions.DryRun, "dry-run", false, "run in dry run mode")
	terraformCmd.AddCommand(terraformPullCmd)
}
