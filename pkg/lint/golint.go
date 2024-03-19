package lint

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/elhub/gh-dxp/pkg/utils"
)

type GoLint struct{}

var golintRegex = regexp.MustCompile(`^([^:]*):([0-9]*):([0-9]*)?:? (.*) \((.*)\)$`)

func (GoLint) Exec(exec *utils.Executor) ([]LinterOutput, error) {
	fmt.Print("Running golint... ")
	var outputs []LinterOutput

	// Run the linter
	out, err := exec.Run("golangci-lint", "run", "./...")
	if err != nil {
		fmt.Printf("Return error: %s\n", err)
	}

	// Go through line by line and parse into lintOuput
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		output, parseErr := GoLintParser(line)
		if parseErr == nil {
			outputs = append(outputs, output)
		}
	}

	fmt.Println("done.")
	return outputs, err
}

func GoLintParser(inputLine string) (LinterOutput, error) {
	matchExtract := golintRegex.FindAllStringSubmatch(inputLine, -1)
	if len(matchExtract) > 0 {
		matches := matchExtract[0]
		line, _ := strconv.Atoi(matches[2])
		column, _ := strconv.Atoi(matches[3])
		return LinterOutput{
			Linter:      "golint",
			Path:        matches[1],
			Line:        line,
			Column:      column,
			Description: matches[4],
			Severity:    "error",
			Source:      matches[5],
		}, nil
	}
	return LinterOutput{}, fmt.Errorf("invalid golint line format")
}
