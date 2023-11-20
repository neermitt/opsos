package cmds

import (
	"fmt"

	"github.com/neermitt/opsos/pkg/plugins/terraform/exec"
	"github.com/neermitt/opsos/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	terraformGenerateVarfileOptions exec.TerraformGenerateVarfileOptions
)

// terraformGenerateVarfileCmd generates varfile for a terraform configuration
var terraformGenerateVarfileCmd = &cobra.Command{
	Use:   "varfile <stack> <component>",
	Short: "Execute 'terraform generate varfile' commands",
	Long:  `This command generates the backend config for a terraform component: opsos terraform generate varfile <stack> <component>`,
	Args:  cobra.MinimumNArgs(2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if terraformGenerateBackendOptions.Format != "" &&
			!utils.StringInSlice(terraformGenerateBackendOptions.Format, []string{"hcl", "json"}) {
			return fmt.Errorf("invalid `format` value `%s`, should be one of `hcl` && `json`", terraformGenerateBackendOptions.Format)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]
		component := args[1]
		return exec.ExecuteTerraformGenerateVarfile(cmd.Context(), stackName, component, terraformGenerateVarfileOptions)
	},
}

func init() {
	terraformGenerateVarfileCmd.Flags().BoolVar(&terraformGenerateVarfileOptions.DryRun, "dry-run", false, "run in dry run mode")
	terraformGenerateVarfileCmd.PersistentFlags().StringVar(&terraformGenerateVarfileOptions.Format, "format", "json", "Specify output format: opsos generate backend --format=hcl/json")
	terraformGenerateCmd.AddCommand(terraformGenerateVarfileCmd)
}
