package cmds

import (
	"github.com/neermitt/opsos/pkg/plugins/terraform/exec"
	"github.com/spf13/cobra"
)

var (
	terraformOptions exec.TerraformOptions
)

// terraformPlanCmd prepares the plan file for the terraform component
var terraformPlanCmd = &cobra.Command{
	Use:   "plan <stack> <component>",
	Short: "Execute 'terraform plan' commands",
	Long:  `This command prepares plan file for a terraform component: opsos terraform plan <stack> <component>`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		terraformOptions.Command = "plan"
		terraformOptions.RequiresVarFile = true
		stackName := args[0]
		component := args[1]
		additionalArgs := args[2:]
		return exec.ExecuteTerraform(cmd.Context(), stackName, component, additionalArgs, terraformOptions)
	},
}

func init() {
	terraformPlanCmd.Flags().BoolVar(&terraformOptions.DryRun, "dry-run", false, "run in dry run mode")
	terraformCmd.AddCommand(terraformPlanCmd)
}
