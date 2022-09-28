package config

type Configuration struct {
	BasePath   string     `yaml:"base_path" json:"base_path" mapstructure:"base_path"`
	Components Components `yaml:"components" json:"components" mapstructure:"components"`
	Stacks     Stacks     `yaml:"stacks" json:"stacks" mapstructure:"stacks"`
	Workflows  Workflows  `yaml:"workflows" json:"workflows" mapstructure:"workflows"`
	Logs       Logs       `yaml:"logs" json:"logs" mapstructure:"logs"`
}

type Components struct {
	Terraform Terraform `yaml:"terraform" json:"terraform" mapstructure:"terraform"`
	Helmfile  Helmfile  `yaml:"helmfile" json:"helmfile" mapstructure:"helmfile"`
}

type Terraform struct {
	BasePath                string `yaml:"base_path" json:"base_path" mapstructure:"base_path"`
	ApplyAutoApprove        bool   `yaml:"apply_auto_approve" json:"apply_auto_approve" mapstructure:"apply_auto_approve"`
	DeployRunInit           bool   `yaml:"deploy_run_init" json:"deploy_run_init" mapstructure:"deploy_run_init"`
	InitRunReconfigure      bool   `yaml:"init_run_reconfigure" json:"init_run_reconfigure" mapstructure:"init_run_reconfigure"`
	AutoGenerateBackendFile bool   `yaml:"auto_generate_backend_file" json:"auto_generate_backend_file" mapstructure:"auto_generate_backend_file"`
}

type Helmfile struct {
	BasePath           string            `yaml:"base_path" json:"base_path" mapstructure:"base_path"`
	KubeconfigPath     string            `yaml:"kubeconfig_path" json:"kubeconfig_path" mapstructure:"kubeconfig_path"`
	ClusterNamePattern string            `yaml:"cluster_name_pattern" json:"cluster_name_pattern" mapstructure:"cluster_name_pattern"`
	Envs               map[string]string `yaml:"envs" json:"envs" mapstructure:"envs"`
}

type Stacks struct {
	BasePath      string   `yaml:"base_path" json:"base_path" mapstructure:"base_path" validate:"required"`
	IncludedPaths []string `yaml:"included_paths" json:"included_paths" mapstructure:"included_paths" validate:"required"`
	ExcludedPaths []string `yaml:"excluded_paths" json:"excluded_paths" mapstructure:"excluded_paths"`
	NamePattern   string   `yaml:"name_pattern" json:"name_pattern" mapstructure:"name_pattern"`
}

type Workflows struct {
	BasePath string `yaml:"base_path" json:"base_path" mapstructure:"base_path"`
}

type Logs struct {
	Level  string `yaml:"verbose" json:"verbose" mapstructure:"verbose"`
	Colors bool   `yaml:"colors" json:"colors" mapstructure:"colors"`
}
