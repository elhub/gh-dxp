package main

import (
	"os"

	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/cmd"
	"github.com/elhub/gh-dxp/pkg/config"
)

var (
	version = "0.1.0"
)

func main() {
	// Use ReadConfig to read .devxp yaml file and print out content
	settings, err := config.ReadConfig(".devxp")
	if err != nil {
		panic(err)
	}

	if cmdErr := cmd.Execute(settings, version); cmdErr != nil {
		log.WithError(err).Error("Command failed")
		os.Exit(1)
	}
}
