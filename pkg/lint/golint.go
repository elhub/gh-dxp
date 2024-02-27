package lint

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/elhub/gh-devxp/pkg/utils"
)

type GoLint struct{}

func (GoLint) Exec() []lintOutput {
	fmt.Println("Running golint")
	var outputs []lintOutput

	// Run the linter
	out, err := utils.Exec("golangci-lint", "run", "./...")
	if err != nil {
		fmt.Printf("Return error: %s\n", err)
	}

	// Go through line by line and parse into lintOuput
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		tokens := strings.Split(line, ":")
		if len(tokens) >= 5 {
			line, _ := strconv.Atoi(tokens[1])
			character, _ := strconv.Atoi(tokens[2])
			description := strings.Join(tokens[4:], " ")
			output := lintOutput{
				linter:      "golint",
				path:        tokens[0],
				line:        line,
				character:   character,
				code:        tokens[3],
				description: description,
			}
			outputs = append(outputs, output)
		}
	}

	return outputs
}
