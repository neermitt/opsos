package components

import (
	"path"

	"github.com/neermitt/opsos/pkg/config"
)

func GetWorkingDirectory(conf *config.Configuration, componentType string, component string) string {
	componentTypeBasePath := conf.Components.Terraform.BasePath
	if componentType == "helmfile" {
		componentTypeBasePath = conf.Components.Helmfile.BasePath
	}
	return path.Join(conf.BasePath, componentTypeBasePath, component)
}
