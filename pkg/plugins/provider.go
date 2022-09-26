package plugins

import (
	"context"

	"github.com/neermitt/opsos/pkg/stack/schema"
)

type ExecutionContext struct {
}

type Provider interface {
	InitStackContext(background context.Context, config schema.StackConfig) context.Context
	ProcessComponentConfig(ctx context.Context, componentConfig schema.ComponentConfig) (map[string]any, error)
}

var providers map[string]Provider

func init() {
	providers = make(map[string]Provider)
}

func RegisterProvider(name string, provider Provider) {
	providers[name] = provider
}

func GetProviders() []string {
	keys := make([]string, 0, len(providers))
	for k := range providers {
		keys = append(keys, k)
	}
	return keys
}

func GetProvider(name string) (Provider, bool) {
	p, ok := providers[name]
	return p, ok
}
