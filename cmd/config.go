package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd configures the opsos CLI
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Execute 'config' commands",
	Long:  `This command shows configuration for CLI`,
}

func init() {
	RootCmd.AddCommand(configCmd)
}
