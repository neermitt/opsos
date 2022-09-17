package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/go-playground/validator"
	"github.com/mitchellh/go-homedir"
	"github.com/neermitt/opsos/pkg/globals"
	"github.com/neermitt/opsos/pkg/logging"
	"github.com/neermitt/opsos/pkg/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// InitConfig finds and merges CLI configurations in the following order: system dir, home dir, current dir, ENV vars, command-line arguments
// https://dev.to/techschoolguru/load-config-from-file-environment-variables-in-golang-with-viper-2j2d
// https://medium.com/@bnprashanth256/reading-configuration-files-and-environment-variables-in-go-golang-c2607f912b63
func InitConfig() (*Configuration, error) {
	// Config is loaded from the following locations (from lower to higher priority):
	// system dir (`/usr/local/etc/opsos` on Linux, `%LOCALAPPDATA%/opsos` on Windows)
	// home dir (~/.opsos)
	// current directory
	// ENV vars
	// Command-line arguments

	logging.Logger.Info("\nSearching, processing and merging opsos CLI configurations (opsos.yaml) in the following order:", zap.Strings("order", []string{"system dir", "home dir", "current dir", "ENV vars", "command-line arguments"}))

	viper.SetEnvPrefix("OPSOS")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yaml")
	viper.SetTypeByDefaultValue(true)

	viper.SetDefault("base_path", "")
	viper.SetDefault("stacks.base_path", "")
	viper.SetDefault("stacks.included_paths", nil)
	viper.SetDefault("stacks.excluded_paths", nil)
	viper.SetDefault("stacks.name_pattern", "")
	viper.SetDefault("components.terraform.base_path", "")
	viper.SetDefault("components.terraform.apply_auto_approve", false)
	viper.SetDefault("components.terraform.deploy_run_init", false)
	viper.SetDefault("components.terraform.auto_generate_backend_file", false)
	viper.SetDefault("components.helmfile.base_path", "")
	viper.SetDefault("components.helmfile.kube_config_path", "")
	viper.SetDefault("components.helmfile.cluster_name_pattern", "")
	viper.SetDefault("workflows.base_path", "")

	// Process config in home dir
	homeDir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	// Process config in the current dir
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
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
			found, err := processConfigFile(cd)
			if err != nil {
				return nil, err
			}
			if found {
				configFound = true
			}
		}
	}

	if !configFound {
		return nil, fmt.Errorf("%s' CLI config files not found in any of the searched paths: system dir, home dir, current dir, ENV vars", globals.ConfigFileName)
	}

	// https://gist.github.com/chazcheadle/45bf85b793dea2b71bd05ebaa3c28644
	// https://sagikazarmark.hu/blog/decoding-custom-formats-with-viper/
	var config Configuration
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	validate := validator.New()
	if err = validate.Struct(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// https://github.com/NCAR/go-figure
// https://github.com/spf13/viper/issues/181
// https://medium.com/@bnprashanth256/reading-configuration-files-and-environment-variables-in-go-golang-c2607f912b63
func processConfigFile(configDir string) (bool, error) {
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

	err = viper.MergeConfig(reader)
	if err != nil {
		return false, err
	}

	logger.Info("Processed CLI config")

	return true, nil
}
