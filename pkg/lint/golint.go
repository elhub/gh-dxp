package lint

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/elhub/gh-devxp/pkg/utils"
)

type GoLint struct{}

func (GoLint) Exec() ([]LintOutput, error) {
	fmt.Print("Running golint... ")
	var outputs []LintOutput

	// Run the linter
	out, err := utils.Exec("golangci-lint", "run", "./...")
	if err != nil {
		fmt.Printf("Return error: %s\n", err)
	}

	// Go through line by line and parse into lintOuput
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		tokens := strings.Split(line, ":")
		if len(tokens) >= 3 {
			line, _ := strconv.Atoi(tokens[1])
			var character int
			var description string
			if len(tokens) >= 4 {
				character, _ = strconv.Atoi(tokens[2])
				description = strings.Join(tokens[3:], " ")
			} else {
				character = 0
				description = strings.Join(tokens[2:], " ")
			}
			output := LintOutput{
				Linter:      "golint",
				Path:        tokens[0],
				Line:        line,
				Character:   character,
				Description: description,
			}
			outputs = append(outputs, output)
		}
	}

	fmt.Println("done.")
	return outputs, err
}
