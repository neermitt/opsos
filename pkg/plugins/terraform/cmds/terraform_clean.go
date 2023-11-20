package cmds

import (
	"github.com/neermitt/opsos/pkg/plugins/terraform/exec"
	"github.com/spf13/cobra"
)

var (
	terraformCleanOptions exec.TerraformCleanOptions
)

// terraformCleanCmd cleans all temporary terraform files
var terraformCleanCmd = &cobra.Command{
	Use:   "clean <stack> <component>",
	Short: "Execute 'terraform clean' commands",
	Long:  `This command cleans all temporary files for a terraform component: opsos terraform clean <stack> <component>`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]
		component := args[1]
		return exec.ExecuteTerraformClean(cmd.Context(), stackName, component, terraformCleanOptions)
	},
}

func init() {
	terraformCleanCmd.Flags().BoolVar(&terraformCleanOptions.ClearDataDir, "clean-data-dir", false, "clean data dir, if TF_DATA_DIR is specified")
	terraformCmd.AddCommand(terraformCleanCmd)
}
