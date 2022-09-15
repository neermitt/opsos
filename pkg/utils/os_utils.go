package utils

import (
	"fmt"
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/neermitt/opsos/pkg/globals"
)

// PrintErrorToStdErrorAndExit prints errors to std.Error and exits with an error code
func PrintErrorToStdErrorAndExit(err error) {
	if err != nil {
		PrintErrorToStdError(err)
		os.Exit(1)
	}
}

// PrintErrorToStdError prints errors to std.Error
func PrintErrorToStdError(err error) {
	if err != nil {
		c := color.New(color.FgRed)
		_, err2 := c.Fprintln(color.Error, err.Error()+"\n")
		if err2 != nil {
			fmt.Println("Error sending the error message to std.Error:")
			PrintError(err2)
			fmt.Println("Original error message:")
			PrintError(err)
		}
	}
}

// PrintError prints errors to std.Output
func PrintError(err error) {
	if err != nil {
		color.Red("%s\n", err)
	}
}

func GetSystemDir() string {
	// https://pureinfotech.com/list-environment-variables-windows-10/
	// https://docs.microsoft.com/en-us/windows/deployment/usmt/usmt-recognized-environment-variables
	// https://softwareengineering.stackexchange.com/questions/299869/where-is-the-appropriate-place-to-put-application-configuration-files-for-each-p
	// https://stackoverflow.com/questions/37946282/why-does-appdata-in-windows-7-seemingly-points-to-wrong-folder
	if runtime.GOOS == "windows" {
		appDataDir := os.Getenv(globals.WindowsAppDataEnvVar)
		if len(appDataDir) > 0 {
			return appDataDir
		}
	} else {
		return globals.SystemDirConfigFilePath
	}

	return ""
}
