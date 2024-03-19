package lint

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/elhub/gh-devxp/pkg/utils"
)

type YamlLint struct{}

var yamlRegex = regexp.MustCompile(`^([^:]*):([0-9]*):([0-9]*): \[(.*)\] (.*) \((.*)\)$`)

func (YamlLint) Exec() ([]LintOutput, error) {
	fmt.Print("Running yamllint... ")
	var outputs []LintOutput

	// Run the linter
	out, err := utils.Exec("yamllint", "-f", "parsable", ".")
	if err != nil {
		fmt.Printf("Return error: %s\n", err)
	}

	// Go through line by line and parse into lintOuput
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		output, err := YamlLintParser(line)
		if err == nil {
			outputs = append(outputs, output)
		}
	}

	fmt.Println("done.")
	return outputs, err
}

func YamlLintParser(inputLine string) (LintOutput, error) {
	matchExtract := yamlRegex.FindAllStringSubmatch(inputLine, -1)
	if len(matchExtract) > 0 {
		matches := matchExtract[0]
		line, _ := strconv.Atoi(matches[2])
		column, _ := strconv.Atoi(matches[3])
		return LintOutput{
			Linter:      "yamllint",
			Path:        matches[1],
			Line:        line,
			Column:      column,
			Description: matches[5],
			Severity:    matches[4],
			Source:      matches[6],
		}, nil
	}
	return LintOutput{}, fmt.Errorf("invalid yamllint line format")
}
