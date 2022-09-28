package plugins

import (
	"context"

	"github.com/neermitt/opsos/pkg/config"
)

type KubeConfigProvider interface {
	ExportKubeConfig(ctx context.Context, clusterName string, kubeConfigPath string) error
}

type KubeConfigProviderFactory func(conf *config.Configuration) KubeConfigProvider

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
