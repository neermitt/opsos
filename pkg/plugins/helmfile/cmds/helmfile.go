package cmds

import (
	"github.com/spf13/cobra"

	"github.com/neermitt/opsos/pkg/plugins/helmfile/exec"
	_ "github.com/neermitt/opsos/pkg/plugins/kind"
)

var (
	helmfileExecOptions exec.HelmfileExecOptions
)

// describeCmd describes configuration for stacks and components
var helmfileCmd = &cobra.Command{
	Use:   "helmfile",
	Short: "Execute 'helmfile' commands",
	Long:  `This command runs helmfile commands`,
}

func InitCommands(parentCmd *cobra.Command) {
	parentCmd.AddCommand(helmfileCmd)
}
