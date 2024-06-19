// Package owner provides the functionality to get the codeowners of a given path
package owner

import (
	"errors"
	"os"

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

	readmeFile := rootDir + "README.md"

	if utils.FileExists(readmeFile) {
		return readmeFile, nil
	}

	// Get the first file in the root directory
	files, err := utils.ListFilesInDirectory(exe, rootDir)
	if err != nil {
		return "", err
	}
	if len(files) > 0 {
		return rootDir + files[0], nil
	}
	return "", errors.New("no files found in the root directory")
}
