package main

import (
	"github.com/michaeloa/gh-devxp/pkg/cmd"
)

var (
	version = "0.1.0"
)

func main() {
	cmd.Execute(version)
}
