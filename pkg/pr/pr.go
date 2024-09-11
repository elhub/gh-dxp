// Package pr contains the functions and types for the pull request command.
package pr

import (
	"strings"

	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/pkg/errors"
)

// CheckForExistingPR checks if a PR already exists for the current branch.
func CheckForExistingPR(exe utils.Executor, branchID string) (string, error) {
	stdOut, err := exe.GH("pr", "list", "-H", branchID, "--json", "number", "--jq", ".[].number")

	if err != nil {
		log.Debug("Error: " + err.Error())
		return "", errors.New("Failed to find existing PR")
	}

	number := strings.Trim(stdOut.String(), "\n")

	return number, nil
}

// GetPRTitle gets the title of the current PR.
func GetPRTitle(exe utils.Executor) (string, error) {
	stdOut, err := exe.GH("pr", "view", "--json", "title", "--jq", ".title")

	if err != nil {
		return "", errors.New("Error getting PR title")
	}

	title := strings.Trim(stdOut.String(), "\n")

	return title, nil
}
