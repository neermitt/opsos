package cmd

import (
	"github.com/spf13/cobra"
)

// componentCmd describes component commands
var componentCmd = &cobra.Command{
	Use:   "component",
	Short: "Execute 'component' commands",
	Long:  `This command runs component commands`,
}

func init() {
	RootCmd.AddCommand(componentCmd)
}
