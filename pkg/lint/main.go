package lint

import (
	"context"
	"fmt"

	"github.com/elhub/gh-devxp/pkg/config"
)

func Run(ctx context.Context, settings *config.Settings) error {
	linters := map[string]Linter{
		"golint": GoLint{},
	}

	// iterate over settings.Linters and run each one
	var outputs []lintOutput
	for _, lintEntry := range settings.Lint.Linters {
		if linter, ok := linters[lintEntry.Name]; ok {
			output := linter.Exec()
			outputs = append(outputs, output...)
		} else {
			fmt.Printf("Linter %s not found\n", lintEntry.Name)
		}
	}

	// print the outputs
	for _, output := range outputs {
		fmt.Printf("%s:%d:%d: %s: %s\n", output.path, output.line, output.character, output.code, output.description)
	}

	return nil

}
