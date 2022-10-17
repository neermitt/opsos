package cmd

import (
	"fmt"
	"github.com/neermitt/opsos/internal/exec"
	"github.com/neermitt/opsos/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	terraformGenerateBackendOptions exec.TerraformGenerateBackendOptions
)

// terraformGenerateBackendCmd generates backend config for a terraform configuration
var terraformGenerateBackendCmd = &cobra.Command{
	Use:   "backend <stack> <component>",
	Short: "Execute 'terraform generate backend' commands",
	Long:  `This command generates the backend config for a terraform component: opsos terraform generate backend <stack> <component>`,
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
		return exec.ExecuteTerraformGenerateBackend(cmd.Context(), stackName, component, terraformGenerateBackendOptions)
	},
}

func init() {
	terraformGenerateBackendCmd.Flags().BoolVar(&terraformGenerateBackendOptions.DryRun, "dry-run", false, "run in dry run mode")
	terraformGenerateBackendCmd.PersistentFlags().StringVar(&terraformGenerateBackendOptions.Format, "format", "json", "Specify output format: opsos generate backend --format=hcl/json ('json' is default)")
	terraformGenerateCmd.AddCommand(terraformGenerateBackendCmd)
}
