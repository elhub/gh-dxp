package owner

import (
	"fmt"
	"os"

	"github.com/hmarr/codeowners"
)

func Execute(path string) error {
	codeownersFile, err := os.Open(".github/CODEOWNERS")
	if err != nil {
		return err
	}

	ruleset, err := codeowners.ParseFile(codeownersFile)
	if err != nil {
		return err
	}

	rule, err := ruleset.Match(path)
	if err != nil {
		return err
	}

	fmt.Printf("Owners: %v\n", rule.Owners)
	return nil
}
