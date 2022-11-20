package v1

import (
	"github.com/neermitt/opsos/api/common"
)

type ComponentSource struct {
	Uri           string   `yaml:"uri" json:"uri"`
	Version       string   `yaml:"version" json:"version"`
	IncludedPaths []string `yaml:"included_paths,omitempty" json:"included_paths,omitempty"`
	ExcludedPaths []string `yaml:"excluded_paths,omitempty" json:"excluded_paths,omitempty"`
}

type ComponentMixins struct {
	Uri      string `yaml:"uri" json:"uri"`
	Version  string `yaml:"version" json:"version"`
	Filename string `yaml:"filename" json:"filename"`
}

type ComponentSpec struct {
	Source ComponentSource   `yaml:"source" json:"source"`
	Mixins []ComponentMixins `yaml:"mixins,omitempty" json:"mixins,omitempty"`
}

type Component struct {
	common.Object `yaml:",inline" json:",inline"`
	Spec          ComponentSpec `yaml:"spec" json:"spec"`
}
