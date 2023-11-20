package cmds

import (
	"github.com/spf13/cobra"
)

// terraformCmd describes configuration for stacks and components
var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Execute 'terraform' commands",
	Long:  `This command runs terraform commands`,
}

func InitCommands(parentCmd *cobra.Command) {
	parentCmd.AddCommand(terraformCmd)
}
