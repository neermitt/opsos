package cmds

import (
	"github.com/neermitt/opsos/pkg/plugins/terraform/exec"
	"github.com/spf13/cobra"
)

// terraformRefreshCmd applies the terraform component
var terraformRefreshCmd = &cobra.Command{
	Use:   "refresh <stack> <component>",
	Short: "Execute 'terraform refresh' commands",
	Long:  `This command refresh a terraform component: opsos terraform refresh <stack> <component>`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		terraformOptions.Command = "refresh"
		terraformOptions.RequiresVarFile = true
		terraformOptions.CleanPlanFileOnCompletion = true
		stackName := args[0]
		component := args[1]
		additionalArgs := args[2:]
		return exec.ExecuteTerraform(cmd.Context(), stackName, component, additionalArgs, terraformOptions)
	},
}

func init() {
	terraformRefreshCmd.Flags().BoolVar(&terraformOptions.DryRun, "dry-run", false, "run in dry run mode")
	terraformCmd.AddCommand(terraformRefreshCmd)
}
