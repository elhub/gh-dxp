package owner

import (
	"fmt"
	"os"

	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/hmarr/codeowners"
)

func Execute(path string, exe utils.Executor) error {
	gitRoot, err := utils.GetGitRootDirectory(exe)
	if err != nil {
		return err
	}

	codeownersFile, err := os.Open(gitRoot + "/.github/CODEOWNERS")
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
