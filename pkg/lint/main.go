package lint

import (
	"context"
	"fmt"

	"github.com/elhub/gh-devxp/pkg/config"
)

var Linters = map[string]Linter{
	"golint": GoLint{},
}

func Run(ctx context.Context, settings *config.Settings) error {

	// iterate over settings.Linters and run each one
	var outputs []LintOutput
	for _, lintEntry := range settings.Lint.Linters {
		if linter, ok := Linters[lintEntry.Name]; ok {
			output, err := linter.Exec()
			outputs = append(outputs, output...)

			if err != nil {
				fmt.Printf("%s returned %d errors\n", lintEntry.Name, len(outputs))
			}

		} else {
			fmt.Printf("Linter %s not found\n", lintEntry.Name)
		}

	}

	// print the outputs
	for _, output := range outputs {
		fmt.Printf("%s:%d:%d: %s: %s\n", output.Path, output.Line, output.Character, output.Code, output.Description)
	}

	return nil

}
