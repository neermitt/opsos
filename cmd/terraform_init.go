package cmd

import (
	"github.com/neermitt/opsos/internal/exec"
	"github.com/spf13/cobra"
)

// terraformInitCmd initializes the terraform component
var terraformInitCmd = &cobra.Command{
	Use:   "init <stack> <component>",
	Short: "Execute 'terraform init' commands",
	Long:  `This command inits a terraform component: opsos terraform init <stack> <component>`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		terraformOptions.Command = "init"
		terraformOptions.SkipInit = true
		terraformOptions.SkipWorkspace = true
		stackName := args[0]
		component := args[1]
		additionalArgs := args[2:]
		return exec.ExecuteTerraform(cmd.Context(), stackName, component, additionalArgs, terraformOptions)
	},
}

func init() {
	terraformInitCmd.Flags().BoolVar(&terraformOptions.DryRun, "dry-run", false, "run in dry run mode")
	terraformCmd.AddCommand(terraformInitCmd)
}
