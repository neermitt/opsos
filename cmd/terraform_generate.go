package cmd

import (
	"github.com/spf13/cobra"
)

// terraformGenerateCmd generates terraform configurations
var terraformGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Execute 'terraform generate' commands",
	Long:  `This command generates configurations for terraform components`,
}

func init() {
	terraformCmd.AddCommand(terraformGenerateCmd)
}
