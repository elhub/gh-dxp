// Package owner provides the functionality to get the codeowners of a given path
package owner

import (
	"os"

	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/hmarr/codeowners"
)

func Execute(path string, exe utils.Executor) ([]string, error) {
	gitRoot, err := utils.GetGitRootDirectory(exe)
	if err != nil {
		return nil, err
	}

	codeownersFile, err := os.Open(gitRoot + "/.github/CODEOWNERS")
	if err != nil {
		return nil, err
	}

	ruleset, err := codeowners.ParseFile(codeownersFile)
	if err != nil {
		return nil, err
	}

	rule, err := ruleset.Match(path)
	if err != nil {
		return nil, err
	}

	// Convert rule.Owners to []string
	owners := make([]string, len(rule.Owners))
	for i, owner := range rule.Owners {
		owners[i] = owner.String()
	}

	// Return the owners of the file
	return owners, nil
}
