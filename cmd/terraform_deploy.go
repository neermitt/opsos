package cmd

import (
	"github.com/neermitt/opsos/internal/exec"
	"github.com/spf13/cobra"
)

// terraformDeployCmd applies the terraform component with auto approve
var terraformDeployCmd = &cobra.Command{
	Use:   "deploy <stack> <component>",
	Short: "Execute 'terraform deploy' command",
	Long:  `This command apply a terraform component with auto approve: opsos terraform deploy <stack> <component>`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		terraformOptions.RequiresVarFile = true
		terraformOptions.AutoApprove = true
		terraformOptions.CleanPlanFileOnCompletion = true
		stackName := args[0]
		component := args[1]
		additionalArgs := args[2:]
		return exec.ExecuteTerraformApply(cmd.Context(), stackName, component, additionalArgs, terraformOptions)
	},
}

func init() {
	terraformDeployCmd.Flags().BoolVar(&terraformOptions.DryRun, "dry-run", false, "run in dry run mode")
	terraformDeployCmd.Flags().BoolVar(&terraformOptions.UsePlan, "use-plan", false, "use existing plan file")
	terraformCmd.AddCommand(terraformDeployCmd)
}
