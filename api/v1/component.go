package v1

import (
	"fmt"
	"io"

	"github.com/neermitt/opsos/api/common"
	"github.com/neermitt/opsos/pkg/utils"
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

func ReadComponent(r io.Reader) (*Component, error) {
	var component Component

	err := utils.DecodeYaml(r, &component)
	if err != nil {
		return nil, err
	}
	err = validateComponent(component)
	if err != nil {
		return nil, err
	}
	return &component, err
}

func validateComponent(component Component) error {
	if component.ApiVersion != "opsos/v1" || component.Kind != "Component" {
		return fmt.Errorf("no resource found of type %s/%s", component.ApiVersion, component.Kind)
	}
	return nil
}
