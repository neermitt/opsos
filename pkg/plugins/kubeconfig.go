package plugins

import (
	"context"
	"errors"
	"fmt"
	"log"

	v1 "github.com/neermitt/opsos/api/v1"
	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/stack"
)

type KubeConfigProvider interface {
	ExportKubeConfig(ctx context.Context, stk *stack.Stack, providerStackSettings map[string]any, kubeConfigPath string) error
}

type KubeConfigProviderFactory func(conf *v1.ConfigSpec) (KubeConfigProvider, error)

var kubeConfigProviderFactoryMap map[string]KubeConfigProviderFactory

func init() {
	kubeConfigProviderFactoryMap = make(map[string]KubeConfigProviderFactory)
}

func RegisterKubeConfigProvider(name string, factory KubeConfigProviderFactory) {
	kubeConfigProviderFactoryMap[name] = factory
}

func GetKubeConfigProvider(ctx context.Context, name string) (KubeConfigProvider, error) {
	pf, ok := kubeConfigProviderFactoryMap[name]
	if !ok {
		return nil, fmt.Errorf("kubeconfig provider `%s` not found", name)
	}
	return pf(config.GetConfig(ctx))
}

func GetKubeConfig(ctx context.Context, k8sSettings map[string]any, stk *stack.Stack, componentName string, kubeConfigPath string) error {
	log.Printf("[TRACE] Begin loading kubeconfig for stack %[1]s", stk.Id)
	defer func() {
		log.Printf("[TRACE] End loading kubeconfig for stack %[1]s", stk.Id)
	}()
	provider, found := k8sSettings["provider"].(string)
	if !found {
		msg := fmt.Sprintf("k8s settings for provider %[3]s not found for component %[2]s in stack %[1]s", stk.Id, componentName, provider)
		log.Printf("[ERROR] %s", msg)
		return errors.New(msg)
	}
	log.Printf("[INFO] loading kubeconfig for stack %[1]s using provider %[2]s", stk.Id, provider)
	kubeConfigProvider, err := GetKubeConfigProvider(ctx, provider)
	if err != nil {
		return err
	}
	providerSettings := k8sSettings[provider].(map[string]any)

	log.Printf("[INFO] Writing the kubeconfig to file: %[1]s", kubeConfigPath)

	return kubeConfigProvider.ExportKubeConfig(ctx, stk, providerSettings, kubeConfigPath)
}
