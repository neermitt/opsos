package exec

import (
	"os"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/formatters"
	"github.com/spf13/cobra"
)

// ExecuteDescribeConfig executes `describe config` command
func ExecuteDescribeConfig(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	format, err := flags.GetString("format")
	if err != nil {
		return err
	}

	conf := cmd.Context().Value("config").(*config.Configuration)

	err = formatters.Get(format)(os.Stdout, conf)
	if err != nil {
		return err
	}

	return nil
}
