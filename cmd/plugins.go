package cmd

import (
	helmfileCmds "github.com/neermitt/opsos/pkg/plugins/helmfile/cmds"
	terraformCmds "github.com/neermitt/opsos/pkg/plugins/terraform/cmds"
	"github.com/spf13/cobra"
)

type RegisterCmdFunc func(command *cobra.Command)

var (
	plugins = []RegisterCmdFunc{helmfileCmds.InitCommands, terraformCmds.InitCommands}
)

func init() {
	for _, registerCmdFunc := range plugins {
		registerCmdFunc(RootCmd)
	}
}
