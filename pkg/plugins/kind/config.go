package kind

type Config struct {
	ClusterNamePattern string `yaml:"cluster_name_pattern" json:"cluster_name_pattern" mapstructure:"cluster_name_pattern"`
}
