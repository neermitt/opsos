package cmd

import (
	"github.com/neermitt/opsos/internal/exec"
	"github.com/spf13/cobra"
)

// terraformDestroyCmd destroys the terraform component with auto approve
var terraformDestroyCmd = &cobra.Command{
	Use:   "destroy <stack> <component>",
	Short: "Execute 'terraform destroy' command",
	Long:  `This command destroys a terraform component with auto approve: opsos terraform destroy <stack> <component>`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		terraformOptions.Destroy = true
		stackName := args[0]
		component := args[1]
		additionalArgs := args[2:]
		return exec.ExecuteTerraformApply(cmd.Context(), stackName, component, additionalArgs, terraformOptions)
	},
}

func init() {
	terraformCmd.AddCommand(terraformDestroyCmd)
}
