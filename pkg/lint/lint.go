package lint

import (
	"fmt"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/utils"
)

func DefaultLinters() map[string]Linter {
	return map[string]Linter{
		"golint":   GoLint{},
		"yamllint": YamlLint{},
	}
}

func Run(exe utils.Executor, settings *config.Settings, linters map[string]Linter) error {
	// iterate over settings.Linters and run each one
	var outputs []LinterOutput
	for _, lintEntry := range settings.Lint.Linters {
		if linter, ok := linters[lintEntry.Name]; ok {
			output, err := linter.Run(exe)
			outputs = append(outputs, output...)

			if err != nil {
				fmt.Printf("%s returned %d errors\n", lintEntry.Name, len(outputs))
			}
		} else {
			fmt.Printf(" Linter %s not found\n", lintEntry.Name)
		}
	}

	// print the outputs
	for _, output := range outputs {
		fmt.Printf("%s:%d:%d: %s: %s\n", output.Path, output.Line, output.Column, output.Description, output.Severity)
	}

	return nil
}
