package cmds

import (
	"github.com/neermitt/opsos/pkg/plugins/terraform/exec"
	"github.com/spf13/cobra"
)

// terraformApplyCmd applies the terraform component
var terraformApplyCmd = &cobra.Command{
	Use:   "apply <stack> <component>",
	Short: "Execute 'terraform apply' commands",
	Long:  `This command apply a terraform component: opsos terraform apply <stack> <component>`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		terraformOptions.Command = "apply"
		terraformOptions.RequiresVarFile = true
		terraformOptions.CleanPlanFileOnCompletion = true
		stackName := args[0]
		component := args[1]
		additionalArgs := args[2:]
		return exec.ExecuteTerraform(cmd.Context(), stackName, component, additionalArgs, terraformOptions)
	},
}

func init() {
	terraformApplyCmd.Flags().BoolVar(&terraformOptions.DryRun, "dry-run", false, "run in dry run mode")
	terraformApplyCmd.Flags().BoolVar(&terraformOptions.UsePlan, "use-plan", false, "use existing plan file")
	terraformCmd.AddCommand(terraformApplyCmd)
}
