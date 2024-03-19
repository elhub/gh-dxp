package main

import (
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
	cmd.Execute(settings, version)
}
