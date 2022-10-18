package cmd

import (
	"github.com/neermitt/opsos/internal/exec"
	"github.com/spf13/cobra"
)

// terraformImportCmd imports the terraform component
var terraformImportCmd = &cobra.Command{
	Use:   "import <stack> <component> ADDR ID",
	Short: "Execute 'terraform import' command",
	Long:  `This command imports a terraform component: opsos terraform import <stack> <component> ADDR ID`,
	Args:  cobra.MinimumNArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		terraformOptions.Destroy = true
		stackName := args[0]
		component := args[1]
		additionalArgs := args[2:]
		return exec.ExecuteTerraformImport(cmd.Context(), stackName, component, additionalArgs, terraformOptions)
	},
}

func init() {
	terraformCmd.AddCommand(terraformImportCmd)
}
