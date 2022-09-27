package plugins

import (
	"github.com/spf13/cobra"
)

type CommandProvider interface {
	RegisterCommands(cmd *cobra.Command)
}

var cmdProvider map[string]CommandProvider

func init() {
	cmdProvider = make(map[string]CommandProvider)
}

func RegisterCmdProvider(name string, provider CommandProvider) {
	cmdProvider[name] = provider
}

func GetCmdProviders() []string {
	keys := make([]string, 0, len(providers))
	for k := range providers {
		keys = append(keys, k)
	}
	return keys
}

func GetCmdProvider(name string) (CommandProvider, bool) {
	p, ok := cmdProvider[name]
	return p, ok
}
