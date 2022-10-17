package exec

import (
	"path"
	"path/filepath"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/stack"
)

func getComponentWorkingDirectory(conf *config.Configuration, componentType string, componentInfo stack.ConfigWithMetadata) (string, string, error) {
	componentTypeBasePath := conf.Components.Terraform.BasePath
	if componentType == "helmfile" {
		componentTypeBasePath = conf.Components.Helmfile.BasePath
	}
	workingDir := path.Join(conf.BasePath, componentTypeBasePath, componentInfo.Component)
	abs, err := filepath.Abs(workingDir)
	return workingDir, abs, err
}
