package cmd

import (
	"github.com/neermitt/opsos/pkg/plugins"
	_ "github.com/neermitt/opsos/pkg/plugins/helmfile"
	_ "github.com/neermitt/opsos/pkg/plugins/kind"
)

func init() {
	for _, providerName := range plugins.GetCmdProviders() {
		provider, _ := plugins.GetCmdProvider(providerName)
		provider.RegisterCommands(RootCmd)
	}
}
