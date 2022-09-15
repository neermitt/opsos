package config

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/neermitt/opsos/pkg/globals"
	"github.com/neermitt/opsos/pkg/logging"
	"github.com/neermitt/opsos/pkg/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	Config     Configuration
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

	logging.Logger.Info("\nSearching, processing and merging opsos CLI configurations (opsos.yaml) in the following order:", zap.Strings("order", []string{"system dir", "home dir", "current dir", "ENV vars", "command-line arguments"}))

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetTypeByDefaultValue(true)

	// Process config in home dir
	homeDir, err := homedir.Dir()
	if err != nil {
		return err
	}

	// Process config in the current dir
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	configDirs := []string{utils.GetSystemDir(), path.Join(homeDir, ".opsos"), cwd}

	// Process config from the path in ENV var `ATMOS_CLI_CONFIG_PATH`
	configPathEnv := os.Getenv("OPSOS_CLI_CONFIG_PATH")

	if len(configPathEnv) > 0 {
		logging.Logger.Debug("Found ENV var", zap.String("OPSOS_CLI_CONFIG_PATH", configPathEnv))
		configDirs = append(configDirs, configPathEnv)
	}

	configFound := false

	for _, cd := range configDirs {
		if len(cd) > 0 {
			found, err := processConfigFile(cd, v)
			if err != nil {
				return err
			}
			if found {
				configFound = true
			}
		}
	}

	if !configFound {
		return fmt.Errorf("%s' CLI config files not found in any of the searched paths: system dir, home dir, current dir, ENV vars", globals.ConfigFileName)
	}

	// https://gist.github.com/chazcheadle/45bf85b793dea2b71bd05ebaa3c28644
	// https://sagikazarmark.hu/blog/decoding-custom-formats-with-viper/
	err = v.Unmarshal(&Config)
	if err != nil {
		return err
	}

	return nil
}

// https://github.com/NCAR/go-figure
// https://github.com/spf13/viper/issues/181
// https://medium.com/@bnprashanth256/reading-configuration-files-and-environment-variables-in-go-golang-c2607f912b63
func processConfigFile(configDir string, v *viper.Viper) (bool, error) {
	logger := logging.Logger.With(zap.String("path", configDir))

	configFile := path.Join(configDir, globals.ConfigFileName)
	if !utils.FileExists(configFile) {
		logger.Info("No config file found")
		return false, nil
	}

	logger.Info("Found config file")

	reader, err := os.Open(configFile)
	if err != nil {
		return false, err
	}

	defer func(reader *os.File) {
		err := reader.Close()
		if err != nil {
			logging.Logger.Error("error closing file", zap.Error(err))
		}
	}(reader)

	err = v.MergeConfig(reader)
	if err != nil {
		return false, err
	}

	logger.Info("Processed CLI config")

	return true, nil
}
