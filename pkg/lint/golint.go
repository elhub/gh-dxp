package lint

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/elhub/gh-devxp/pkg/utils"
)

type GoLint struct{}

var golintRegex = regexp.MustCompile(`^([^:]*):([0-9]*):([0-9]*)?:? (.*) \((.*)\)$`)

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
		output, err := GoLintParser(line)
		if err == nil {
			outputs = append(outputs, output)
		}
	}

	fmt.Println("done.")
	return outputs, err
}

func GoLintParser(inputLine string) (LintOutput, error) {
	matchExtract := golintRegex.FindAllStringSubmatch(inputLine, -1)
	if len(matchExtract) > 0 {
		matches := matchExtract[0]
		line, _ := strconv.Atoi(matches[2])
		column, _ := strconv.Atoi(matches[3])
		return LintOutput{
			Linter:      "golint",
			Path:        matches[1],
			Line:        line,
			Column:      column,
			Description: matches[4],
			Severity:    "error",
			Source:      matches[5],
		}, nil
	}
	return LintOutput{}, fmt.Errorf("invalid golint line format")
}
