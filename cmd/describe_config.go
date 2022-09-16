package cmd

import (
	"github.com/neermitt/opsos/internal/exec"
	"github.com/neermitt/opsos/pkg/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// describeComponentCmd describes configuration for components
var describeConfigCmd = &cobra.Command{
	Use:                "config",
	Short:              "Execute 'describe config' command",
	Long:               `This command shows the final (deep-merged) CLI configuration: atmos describe config`,
	FParseErrWhitelist: struct{ UnknownFlags bool }{UnknownFlags: true},
	Run: func(cmd *cobra.Command, args []string) {
		err := exec.ExecuteDescribeConfig(cmd, args)
		if err != nil {
			logging.Logger.Error("Decribe Config Failed", zap.Error(err))
		}
	},
}

func init() {
	describeConfigCmd.DisableFlagParsing = false
	describeConfigCmd.PersistentFlags().StringP("format", "f", "json", "'json' or 'yaml'")

	describeCmd.AddCommand(describeConfigCmd)
}
