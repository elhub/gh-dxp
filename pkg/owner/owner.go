// Package owner provides the functionality to get the codeowners of a given path
package owner

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/hmarr/codeowners"
)

// Execute determines the owner of the specified file based on the CODEOWNERS file given in the .github directory.
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

// GetDefaultFile returns the default file to check for codeowners.
func GetDefaultFile(exe utils.Executor) (string, error) {
	rootDir, err := utils.GetGitRootDirectory(exe)
	if err != nil {
		return "", err
	}

	ownersFile := filepath.Join(rootDir, ".github/CODEOWNERS")

	if utils.FileExists(ownersFile) {
		return ownersFile, nil
	}

	return "", errors.New("could not find CODEOWNERS file in .github directory")
}
