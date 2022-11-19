package cmd

import (
	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/logging"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "opsos",
	Short: "Universal Tool for DevOps and Cloud Automation",
	Long:  `'opsos'' is a universal tool for DevOps and cloud automation used for provisioning, managing and orchestrating workflows across various toolchains`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config load if version command
		if cmd == versionCmd {
			return nil
		}
		// InitConfig finds and merges CLI configurations in the following order:
		// system dir, home dir, current dir, ENV vars, command-line arguments
		// Here we need the custom commands from the config
		conf, err := config.InitConfig()
		if err != nil {
			return err
		}

		cmd.SetContext(config.SetConfig(cmd.Context(), conf))
		logging.InitLogger(*conf)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() error {
	defer logging.PanicHandler()
	return RootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {

}

// https://www.sobyte.net/post/2021-12/create-cli-app-with-cobra/
// https://github.com/spf13/cobra/blob/master/user_guide.md
// https://blog.knoldus.com/create-kubectl-like-cli-with-go-and-cobra/
// https://pkg.go.dev/github.com/c-bata/go-prompt
// https://pkg.go.dev/github.com/spf13/cobra
// https://scene-si.org/2017/04/20/managing-configuration-with-viper/
