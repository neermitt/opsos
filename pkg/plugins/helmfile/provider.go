package helmfile

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/merge"
	"github.com/neermitt/opsos/pkg/plugins"
	"github.com/neermitt/opsos/pkg/stack/schema"
	"github.com/spf13/afero"
)

const (
	providerName = "helmfile"
	stackConfig  = "helmfile-stack-config"

	varsField = "vars"
)

type helmfileComponentInfo struct {
	Component             string `yaml:"-" json:"-"`
	ComponentFolderPrefix string `yaml:"-" json:"-"`

	FinalComponent string         `yaml:"FinalComponent,omitempty" json:"FinalComponent,omitempty"`
	Vars           map[string]any `yaml:"vars" json:"vars"`
	WorkingDir     string         `yaml:"-" json:"-"`
}

type helmfileProvider struct {
	fs       afero.Fs
	basePath string
}

func (h *helmfileProvider) InitStackContext(ctx context.Context, config schema.StackConfig) context.Context {
	componentTypeVars, _ := merge.Merge([]map[string]any{config.Vars, config.ComponentTypeSettings[providerName].Vars})
	return context.WithValue(ctx, stackConfig, map[string]any{"vars": componentTypeVars})
}

func (h *helmfileProvider) ProcessComponentConfig(ctx context.Context, componentName string, componentConfig schema.ComponentConfig) (any, error) {
	baseConfig := ctx.Value(stackConfig).(map[string]any)
	baseVars := baseConfig[varsField].(map[string]any)
	componentVars, err := merge.Merge([]map[string]any{baseVars, componentConfig.Vars})
	if err != nil {
		return nil, err
	}

	info := helmfileComponentInfo{Vars: componentVars}

	component := componentName

	if componentConfig.Component != "" {
		info.FinalComponent = componentConfig.Component
		component = componentConfig.Component
	}

	compoDirInfo, err := h.fs.Stat(component)
	if err != nil {
		return nil, err
	}
	if !compoDirInfo.IsDir() {
		return nil, fmt.Errorf("component directory %s not found in base path %s", component, h.basePath)
	}
	dir, file := path.Split(component)
	if dir == "" {
		info.Component = component
	} else {
		info.Component = file
		info.ComponentFolderPrefix = strings.TrimSuffix(dir, "/")
	}

	info.WorkingDir = path.Join(h.basePath, component)

	return info, nil
}

func init() {
	plugins.RegisterProvider(providerName, func(conf *config.Configuration) plugins.Provider {
		helmfileBasePath := path.Join(conf.BasePath, conf.Components.Helmfile.BasePath)
		helmfileBaseAbsPath, err := filepath.Abs(helmfileBasePath)
		if err != nil {
			return nil
		}

		stackFS := afero.NewBasePathFs(afero.NewOsFs(), helmfileBaseAbsPath)

		return &helmfileProvider{fs: stackFS, basePath: helmfileBaseAbsPath}
	})
}
