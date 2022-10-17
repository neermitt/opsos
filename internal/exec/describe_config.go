package exec

import (
	"github.com/neermitt/opsos/pkg/utils"
	"os"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/spf13/cobra"
)

type DescribeConfigOptions struct {
	Format string
}

// ExecuteDescribeConfig executes `describe config` command
func ExecuteDescribeConfig(cmd *cobra.Command, options DescribeConfigOptions) error {
	conf := config.GetConfig(cmd.Context())

	err := utils.Get(options.Format)(os.Stdout, conf)
	if err != nil {
		return err
	}

	return nil
}
