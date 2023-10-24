package cmds

import (
	"github.com/neermitt/opsos/pkg/plugins/terraform/exec"
	"github.com/spf13/cobra"
)

// terraformWorkspaceCmd runs workspace commands
var terraformWorkspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Execute 'terraform workspace' commands",
	Long:  `This command runs workspaces command for terraform components`,
}

var terraformWorkspaceListCmd = &cobra.Command{
	Use:   "list <stack> <component>",
	Short: "Execute 'terraform workspace list' commands",
	Long:  `This command list workspaces for terraform components`,
	Args:  cobra.MinimumNArgs(2),
	RunE:  handleWorkspaceCmd,
}

var terraformWorkspaceShowCmd = &cobra.Command{
	Use:   "show <stack> <component>",
	Short: "Execute 'terraform workspace show' commands",
	Long:  `This command shows name of current workspace for terraform components`,
	Args:  cobra.MinimumNArgs(2),
	RunE:  handleWorkspaceCmd,
}

func init() {
	terraformCmd.AddCommand(terraformWorkspaceCmd)
	terraformWorkspaceCmd.AddCommand(terraformWorkspaceListCmd)
	terraformWorkspaceCmd.AddCommand(terraformWorkspaceShowCmd)
}

func handleWorkspaceCmd(cmd *cobra.Command, args []string) error {
	terraformOptions.Command = "workspace"
	stackName := args[0]
	component := args[1]
	additionalArgs := append([]string{cmd.CalledAs()}, args[2:]...)
	return exec.ExecuteTerraform(cmd.Context(), stackName, component, additionalArgs, terraformOptions)
}
