package cmds

import (
	"github.com/neermitt/opsos/pkg/plugins/terraform/exec"
	"github.com/spf13/cobra"
)

// terraformShellCmd starts a shell for the terraform component
var terraformShellCmd = &cobra.Command{
	Use:   "shell <stack> <component>",
	Short: "Starts 'terraform shell'",
	Long:  `This command starts a shell for a terraform component: opsos terraform shell <stack> <component>`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		terraformOptions.Command = "shell"
		terraformOptions.RequiresVarFile = true
		stackName := args[0]
		component := args[1]
		additionalArgs := args[2:]
		return exec.ExecuteTerraform(cmd.Context(), stackName, component, additionalArgs, terraformOptions)
	},
}

func init() {
	terraformCmd.AddCommand(terraformShellCmd)
}
