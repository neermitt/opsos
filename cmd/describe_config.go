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
	Use:                "config",
	Short:              "Execute 'describe config' command",
	Long:               `This command shows the final (deep-merged) CLI configuration: atmos describe config`,
	FParseErrWhitelist: struct{ UnknownFlags bool }{UnknownFlags: true},
	RunE: func(cmd *cobra.Command, args []string) error {
		return exec.ExecuteDescribeConfig(cmd, descConfOptions)
	},
}

func init() {
	describeConfigCmd.PersistentFlags().StringVarP(&descConfOptions.Format, "format", "f", "yaml", "'json' or 'yaml'")

	describeCmd.AddCommand(describeConfigCmd)
}
