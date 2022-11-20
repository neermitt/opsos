package helmfile

type Config struct {
	BasePath       string            `yaml:"base_path" json:"base_path" mapstructure:"base_path"`
	KubeconfigPath string            `yaml:"kubeconfig_path" json:"kubeconfig_path" mapstructure:"kubeconfig_path"`
	Envs           map[string]string `yaml:"envs" json:"envs" mapstructure:"envs"`
}
