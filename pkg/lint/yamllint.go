package lint

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/elhub/gh-dxp/pkg/utils"
)

type YamlLint struct{}

var yamlRegex = regexp.MustCompile(`^([^:]*):([0-9]*):([0-9]*): \[(.*)\] (.*) \((.*)\)$`)

func (YamlLint) Run(exe utils.Executor) ([]LinterOutput, error) {
	s := utils.StartSpinner("Running Yamllint...", "Yamllint done.")
	var outputs []LinterOutput

	// Run the linter
	out, err := exe.Command("yamllint", "-f", "parsable", ".")
	if err != nil {
		fmt.Printf("Return error: %s\n", err)
	}

	// Go through line by line and parse into lintOuput
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		output, parseErr := YamlLintParser(line)
		if parseErr == nil {
			outputs = append(outputs, output)
		}
	}

	s.Stop()
	return outputs, err
}

func YamlLintParser(inputLine string) (LinterOutput, error) {
	matchExtract := yamlRegex.FindAllStringSubmatch(inputLine, -1)
	if len(matchExtract) > 0 {
		matches := matchExtract[0]
		line, _ := strconv.Atoi(matches[2])
		column, _ := strconv.Atoi(matches[3])
		return LinterOutput{
			Linter:      "yamllint",
			Path:        matches[1],
			Line:        line,
			Column:      column,
			Description: matches[5],
			Severity:    matches[4],
			Source:      matches[6],
		}, nil
	}
	return LinterOutput{}, fmt.Errorf("invalid yamllint line format")
}
