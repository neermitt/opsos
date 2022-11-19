package cmd

import "github.com/spf13/cobra"

// stackCmd describes stack commands
var stackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Execute 'stack' commands",
	Long:  `This command runs stack commands`,
}

func init() {
	RootCmd.AddCommand(stackCmd)
}
