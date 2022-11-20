package components

import (
	"path"

	v1 "github.com/neermitt/opsos/api/v1"
)

func GetWorkingDirectory(conf *v1.ConfigSpec, componentType string, component string) string {
	if conf.Providers[componentType] == nil {
		return ""
	}
	componentTypeBasePath := conf.Providers[componentType]["base_path"].(string)
	return path.Join(*conf.BasePath, componentTypeBasePath, component)
}
