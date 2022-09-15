package main

import (
	"github.com/neermitt/opsos/cmd"
	"github.com/neermitt/opsos/pkg/utils"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		utils.PrintErrorToStdErrorAndExit(err)
	}
}
