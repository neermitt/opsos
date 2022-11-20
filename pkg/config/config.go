package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator"
	"github.com/mitchellh/go-homedir"
	v1 "github.com/neermitt/opsos/api/v1"
	"github.com/neermitt/opsos/pkg/globals"
	"github.com/neermitt/opsos/pkg/merge"
	"github.com/neermitt/opsos/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	envConfigPath = "OPSOS_CONFIG_PATH"
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

	log.Printf("[INFO] Searching, processing and merging opsos CLI configurations (%s) in the following order: %v",
		globals.ConfigFileName,
		[]string{"system dir", "home dir", "current dir", "ENV vars", "command-line arguments"})

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
	viper.SetDefault("logs.level", "INFO")
	viper.SetDefault("logs.json", false)
	viper.SetDefault("logs.file", nil)

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

	// Process config from the path in ENV
	configPathEnv := os.Getenv(envConfigPath)

	if len(configPathEnv) > 0 {
		log.Printf("[INFO] Found ENV var %s=%s", envConfigPath, configPathEnv)
		configDirs = append(configDirs, configPathEnv)
	}

	configDirs = utils.Unique(configDirs)
	conf, err := ReadAndMergeConfigsFromDirs(configDirs)
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

func ReadAndMergeConfigsFromDirs(dirs []string) (*v1.Config, error) {
	configs := make([]*v1.Config, 0)
	for _, dir := range dirs {
		opsosConfigFileName := filepath.Join(dir, globals.ConfigFileName)
		if utils.FileExists(opsosConfigFileName) {
			log.Printf("[DEBUG] Found config file at %s", dir)
			config, err := readConfigFromFile(opsosConfigFileName)
			if err != nil {
				return nil, errors.Wrapf(err, "Invalid config file %s", opsosConfigFileName)
			}
			configs = append(configs, config)
		}
	}

	return mergeConfigs(configs)
}

func mergeConfigs(configs []*v1.Config) (*v1.Config, error) {
	switch len(configs) {
	case 0:
		return nil, nil
	case 1:
		return configs[0], nil
	}

	log.Print("[DEBUG] Merging multiple configs")
	specs := make([]map[string]any, len(configs))
	for i, config := range configs {
		var err error
		specs[i], err = utils.ToMap(config.Spec)
		if err != nil {
			return nil, err
		}
	}
	mergedSpec, err := merge.Merge(specs)
	if err != nil {
		return nil, err
	}
	targetConfig := configs[len(configs)]
	err = utils.FromMap(mergedSpec, &targetConfig.Spec)
	if err != nil {
		return nil, err
	}
	return targetConfig, nil
}

func readConfigFromFile(filename string) (*v1.Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return readConfig(file)
}

func readConfig(r io.Reader) (*v1.Config, error) {
	var config v1.Config

	err := utils.DecodeYaml(r, &config)
	if err != nil {
		return nil, err
	}
	err = validateConfig(config)
	if err != nil {
		return nil, err
	}
	return &config, err
}

func validateConfig(component v1.Config) error {
	if component.ApiVersion != "opsos/v1" || component.Kind != "Configuration" {
		return fmt.Errorf("no resource found of type %s/%s", component.ApiVersion, component.Kind)
	}
	return nil
}
