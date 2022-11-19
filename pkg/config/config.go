package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/go-playground/validator"
	"github.com/mitchellh/go-homedir"
	v1 "github.com/neermitt/opsos/api/v1"
	"github.com/neermitt/opsos/pkg/globals"
	"github.com/neermitt/opsos/pkg/logging"
	"github.com/neermitt/opsos/pkg/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// InitConfig finds and merges CLI configurations in the following order: system dir, home dir, current dir, ENV vars, command-line arguments
// https://dev.to/techschoolguru/load-config-from-file-environment-variables-in-golang-with-viper-2j2d
// https://medium.com/@bnprashanth256/reading-configuration-files-and-environment-variables-in-go-golang-c2607f912b63
func InitConfig() (*v1.ConfigSpec, error) {
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
	viper.SetDefault("terraform.base_path", "")
	viper.SetDefault("terraform.apply_auto_approve", false)
	viper.SetDefault("terraform.deploy_run_init", false)
	viper.SetDefault("terraform.auto_generate_backend_file", false)
	viper.SetDefault("helmfile.base_path", "")
	viper.SetDefault("helmfile.kube_config_path", "")
	viper.SetDefault("helmfile.cluster_name_pattern", "")
	viper.SetDefault("helmfile.envs", nil)
	viper.SetDefault("kind.cluster_name_pattern", "")

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

	conf, err := v1.ReadAndMergeConfigsFromDirs(configDirs)
	if err != nil {
		return nil, err
	}

	if conf == nil {
		return nil, fmt.Errorf("%s' CLI config files not found in any of the searched paths: system dir, home dir, current dir, ENV vars", globals.ConfigFileName)
	}

	yamlConfig, err := utils.ConvertToYAML(conf.Spec)
	if err != nil {
		return nil, err
	}

	err = viper.MergeConfig(strings.NewReader(yamlConfig))
	if err != nil {
		return nil, err
	}

	var confSpec v1.ConfigSpec
	err = viper.Unmarshal(&confSpec)
	if err != nil {
		return nil, err
	}

	validate := validator.New()
	if err = validate.Struct(&conf.Spec); err != nil {
		return nil, err
	}

	return &confSpec, nil
}
