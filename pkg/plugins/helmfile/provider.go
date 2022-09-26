package helmfile

import (
	"context"

	"github.com/neermitt/opsos/pkg/merge"
	"github.com/neermitt/opsos/pkg/plugins"
	"github.com/neermitt/opsos/pkg/stack/schema"
)

const providerName = "helmfile"

type helmfileProvider struct {
}

func (h *helmfileProvider) InitStackContext(ctx context.Context, config schema.StackConfig) context.Context {
	componentTypeVars, _ := merge.Merge([]map[string]any{config.Vars, config.ComponentTypeSettings[providerName].Vars})
	return context.WithValue(ctx, providerName, map[string]any{"vars": componentTypeVars})
}

func (h *helmfileProvider) ProcessComponentConfig(ctx context.Context, componentConfig schema.ComponentConfig) (map[string]any, error) {
	baseConfig := ctx.Value(providerName).(map[string]any)
	baseVars := baseConfig["vars"].(map[string]any)
	componentVars, err := merge.Merge([]map[string]any{baseVars, componentConfig.Vars})
	if err != nil {
		return nil, err
	}
	return map[string]any{"vars": componentVars}, nil
}

func init() {
	plugins.RegisterProvider(providerName, &helmfileProvider{})
}
