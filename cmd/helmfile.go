package cmd

import (
	"github.com/neermitt/opsos/internal/exec"
	"github.com/spf13/cobra"

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

func init() {
	RootCmd.AddCommand(helmfileCmd)
}
