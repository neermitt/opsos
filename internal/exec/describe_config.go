package exec

import (
	"os"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/formatters"
	"github.com/spf13/cobra"
)

type DescribeConfigOptions struct {
	Format string
}

// ExecuteDescribeConfig executes `describe config` command
func ExecuteDescribeConfig(cmd *cobra.Command, options DescribeConfigOptions) error {
	conf := config.GetConfig(cmd.Context())

	err := formatters.Get(options.Format)(os.Stdout, conf)
	if err != nil {
		return err
	}

	return nil
}
