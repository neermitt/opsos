package kind

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	v1 "github.com/neermitt/opsos/api/v1"
	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/plugins"
	"github.com/neermitt/opsos/pkg/stack"
	"github.com/neermitt/opsos/pkg/utils"
)

func NewKindKubeConfigProvider(_ *v1.ConfigSpec) (plugins.KubeConfigProvider, error) {
	return &kindKubeConfigProvider{}, nil
}

func init() {
	plugins.RegisterKubeConfigProvider("kind", NewKindKubeConfigProvider)
}

type kindKubeConfigProvider struct {
}

func (k *kindKubeConfigProvider) ExportKubeConfig(ctx context.Context, stk *stack.Stack, providerStackSettings map[string]any, kubeConfigPath string) error {
	conf := config.GetConfig(ctx)
	var kindConfig Config
	err := utils.FromMap(conf.Providers[ComponentType], &kindConfig)
	if err != nil {
		return err
	}
	clusterName, err := utils.ProcessTemplate(kindConfig.ClusterNamePattern, stk.Vars)
	if err != nil {
		return err
	}
	err = exportKindKubeConfigRaw(ctx, clusterName, kubeConfigPath)
	if err != nil {
		return err
	}

	controlPlaneIPAddress, err := getKindControlPlaneIPAddress(ctx, clusterName)
	if err != nil {
		return err
	}

	// kubectl config set clusters.kind-"${1}".server https://"${IP}":6443
	err = utils.ExecuteShellCommand(ctx, "kubectl", []string{
		"config",
		"set",
		fmt.Sprintf("clusters.kind-%s.server", clusterName),
		fmt.Sprintf("https://%s:6443", controlPlaneIPAddress)},
		utils.ExecOptions{
			Env: []string{fmt.Sprintf("KUBECONFIG=%s", kubeConfigPath)},
		})

	if err != nil {
		return err
	}

	return nil
}

func exportKindKubeConfigRaw(ctx context.Context, clusterName string, kubeConfigPath string) (err error) {
	f, err := os.OpenFile(kubeConfigPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := f.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	err = utils.ExecuteShellCommand(ctx, "kind", []string{"get", "kubeconfig", "--name", clusterName}, utils.ExecOptions{
		StdOut: f,
	})
	if err != nil {
		return err
	}

	return nil
}

func getKindControlPlaneIPAddress(ctx context.Context, clusterName string) (string, error) {
	var out bytes.Buffer
	err := utils.ExecuteShellCommand(ctx, "docker", []string{
		"inspect",
		"-f",
		"{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}",
		fmt.Sprintf("%s-control-plane", clusterName)},
		utils.ExecOptions{
			StdOut: &out,
		})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}
