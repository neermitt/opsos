package cmd

import (
	"github.com/neermitt/opsos/internal/exec"
	"github.com/spf13/cobra"
)

// terraformApplyCmd applies the terraform component
var terraformApplyCmd = &cobra.Command{
	Use:     "apply <stack> <component>",
	Aliases: []string{"deploy"},
	Short:   "Execute 'terraform apply' commands",
	Long:    `This command apply a terraform component: opsos terraform apply <stack> <component>`,
	Args:    cobra.MinimumNArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
		if cmd.CalledAs() == "deploy" {
			terraformOptions.AutoApprove = true
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]
		component := args[1]
		additionalArgs := args[2:]
		return exec.ExecuteTerraformApply(cmd.Context(), stackName, component, additionalArgs, terraformOptions)
	},
}

func init() {
	terraformApplyCmd.Flags().BoolVar(&terraformOptions.UsePlan, "use-plan", false, "use existing plan file")
	terraformCmd.AddCommand(terraformApplyCmd)
}
