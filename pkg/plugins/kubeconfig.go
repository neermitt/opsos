package plugins

import (
	"context"
	"fmt"

	v1 "github.com/neermitt/opsos/api/v1"
	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/stack"
)

type KubeConfigProvider interface {
	ExportKubeConfig(ctx context.Context, stk *stack.Stack, kubeConfigPath string) error
}

type KubeConfigProviderFactory func(conf *v1.ConfigSpec) KubeConfigProvider

var kubeConfigProviderFactoryMap map[string]KubeConfigProviderFactory

func init() {
	kubeConfigProviderFactoryMap = make(map[string]KubeConfigProviderFactory)
}

func RegisterKubeConfigProvider(name string, factory KubeConfigProviderFactory) {
	kubeConfigProviderFactoryMap[name] = factory
}

func GetKubeConfigProvider(ctx context.Context, name string) (KubeConfigProvider, bool) {
	pf, ok := kubeConfigProviderFactoryMap[name]
	if !ok {
		return nil, ok
	}
	return pf(config.GetConfig(ctx)), ok
}

func GetKubeConfig(ctx context.Context, provider string, stk *stack.Stack, kubeConfigPath string) error {
	kubeConfigProvider, found := GetKubeConfigProvider(ctx, provider)
	if !found {
		return fmt.Errorf("%s kube config provider is not configured", provider)
	}
	return kubeConfigProvider.ExportKubeConfig(ctx, stk, kubeConfigPath)
}
