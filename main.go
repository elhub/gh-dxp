package main

import (
	"fmt"

	"github.com/michaeloa/gh-devxp/pkg/cmd"
	"github.com/michaeloa/gh-devxp/pkg/config"
)

var (
	version = "0.1.0"
)

func main() {
	// Use ReadConfig to read .devxp yaml file and print out content
	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}
	// Print contents of cfg to stdout
	fmt.Printf("%+v\n", cfg)
	if len(cfg.Lint.Linters) > 0 {
		fmt.Printf("%+v\n", cfg.Lint.Linters[0])
	} else {
		fmt.Println("cfg is empty")
	}
	cmd.Execute(version)
}
