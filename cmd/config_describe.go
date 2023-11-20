package cmd

import (
	"github.com/neermitt/opsos/internal/exec"
	"github.com/spf13/cobra"
)

var (
	descConfOptions exec.DescribeConfigOptions
)

// describeComponentCmd describes configuration for components
var describeConfigCmd = &cobra.Command{
	Use:                "describe",
	Short:              "Execute 'config describe' command",
	Long:               `This command shows the final (deep-merged) CLI configuration: opsos config describe`,
	FParseErrWhitelist: struct{ UnknownFlags bool }{UnknownFlags: true},
	RunE: func(cmd *cobra.Command, args []string) error {
		return exec.ExecuteDescribeConfig(cmd, descConfOptions)
	},
}

func init() {
	describeConfigCmd.PersistentFlags().StringVarP(&descConfOptions.Format, "format", "f", "yaml", "'json' or 'yaml'")

	configCmd.AddCommand(describeConfigCmd)
}
