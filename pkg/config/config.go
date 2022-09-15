package config

import (
	"github.com/neermitt/opsos/pkg/logging"
	"go.uber.org/zap"
)

var (
	config     Configuration
	intialized bool = false
)

// InitConfig finds and merges CLI configurations in the following order: system dir, home dir, current dir, ENV vars, command-line arguments
// https://dev.to/techschoolguru/load-config-from-file-environment-variables-in-golang-with-viper-2j2d
// https://medium.com/@bnprashanth256/reading-configuration-files-and-environment-variables-in-go-golang-c2607f912b63
func InitConfig() error {
	// Config is loaded from the following locations (from lower to higher priority):
	// system dir (`/usr/local/etc/opsos` on Linux, `%LOCALAPPDATA%/opsos` on Windows)
	// home dir (~/.opsos)
	// current directory
	// ENV vars
	// Command-line arguments

	if intialized {
		return nil
	}

	err := processLogsConfig()
	if err != nil {
		return err
	}

	logging.Logger.Info("\nSearching, processing and merging opsos CLI configurations (opsos.yaml) in the following order:", zap.Strings("order", []string{"system dir", "home dir", "current dir", "ENV vars", "command-line arguments"}))

	return nil
}

func processLogsConfig() error {
	logging.InitLogger()

	return nil
}
