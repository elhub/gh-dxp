package lint

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/elhub/gh-dxp/pkg/utils"
)

type Detekt struct{}

var detektRegex = regexp.MustCompile(`^\s*(.*)\s*-\s*\[(.*)\]\s*at\s*(.*):([0-9]*):([0-9]*)$`)

func (Detekt) Run(exe utils.Executor) ([]LinterOutput, error) {
	s := utils.StartSpinner("Running Detekt...", "Detekt done.")
	var outputs []LinterOutput

	// Run the linter
	lintString, filesErr := GetFiles(",", ".kt")
	if filesErr != nil {
		fmt.Printf("Return error: %s\n", filesErr)
	}
	_, err := exe.Command("detekt", "-i", lintString, "-r", "md:detekt.out")
	if err != nil {
		fmt.Printf("Return error: %s\n", err)
	}

	// Read output file
	data, readErr := os.ReadFile("detekt.out")
	if readErr != nil {
		return nil, readErr
	}

	// Go through output file line by line and parse into lintOuput
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		output, parseErr := DetektParser(line)
		if parseErr == nil {
			outputs = append(outputs, output)
		}
	}

	// Clean up
	_, err = exe.Command("rm", "detekt.out")
	if err != nil {
		fmt.Printf("Return error: %s\n", err)
	}

	s.Stop()
	return outputs, err
}

func DetektParser(inputLine string) (LinterOutput, error) {
	matchExtract := detektRegex.FindAllStringSubmatch(inputLine, -1)
	if len(matchExtract) > 0 {
		matches := matchExtract[0]
		line, _ := strconv.Atoi(matches[4])
		column, _ := strconv.Atoi(matches[5])
		return LinterOutput{
			Linter:      "detekt",
			Path:        matches[3],
			Line:        line,
			Column:      column,
			Description: matches[1] + "in " + matches[2],
			Severity:    "error",
			Source:      "detekt",
		}, nil
	}
	return LinterOutput{}, fmt.Errorf("invalid detekt line format")
}
