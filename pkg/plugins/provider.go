package plugins

import (
	"context"

	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/stack/schema"
)

type ExecutionContext struct {
}

type Provider interface {
	InitStackContext(background context.Context, config schema.StackConfig) context.Context
	ProcessComponentConfig(ctx context.Context, componentName string, componentConfig schema.ComponentConfig) (any, error)
}

type ProviderFactory func(conf *config.Configuration) Provider

var providers map[string]ProviderFactory

func init() {
	providers = make(map[string]ProviderFactory)
}

func RegisterProvider(name string, provider ProviderFactory) {
	providers[name] = provider
}

func GetProviders() []string {
	keys := make([]string, 0, len(providers))
	for k := range providers {
		keys = append(keys, k)
	}
	return keys
}

func GetProvider(ctx context.Context, name string) (Provider, bool) {
	pf, ok := providers[name]
	return pf(config.GetConfig(ctx)), ok
}
