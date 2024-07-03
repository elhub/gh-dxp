// Package main contains the main function for the gh-dxp extension.
package main

import (
	"os"

	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/cmd"
	"github.com/elhub/gh-dxp/pkg/config"
)

var (
	version = "locally compiled development build"
)

func main() {
	// If there is a config file in ~/.local/devxp/config.yml, use that.
	// Otherwise, use the default .devxp file in the current directory.
	// If neither is found, use default settings.
	settings := config.DefaultSettings()
	if userSettings, err := config.ReadConfig("~/.local/devxp/config.yml"); err == nil {
		settings = config.MergeSettings(settings, userSettings)
	} else if localSettings, err := config.ReadConfig(".devxp"); err == nil {
		settings = config.MergeSettings(settings, localSettings)
	}

	if cmdErr := cmd.Execute(settings, version); cmdErr != nil {
		log.WithError(cmdErr).Error("Command failed")
		os.Exit(1)
	}
}
