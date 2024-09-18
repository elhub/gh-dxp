package pr

import (
	"strings"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/pkg/errors"
)

// ExecuteUpdate updates a pull request, depending on its current state.
func ExecuteUpdate(exe utils.Executor, settings *config.Settings, options *UpdateOptions) error {
	pr := PullRequest{}

	// Get branchID
	currentBranch, errBranch := exe.Command("git", "branch", "--show-current")
	if errBranch != nil {
		return errBranch
	}
	pr.branchID = strings.Trim(currentBranch, "\n")

	opts := &CreateOptions{
		TestRun: options.TestRun,
		NoLint:  options.NoLint,
		NoUnit:  options.NoUnit,
	}

	// Check if PR exists on branch
	prID, errCheck := CheckForExistingPR(exe, pr.branchID)
	if errCheck != nil {
		return errCheck
	}

	pr, err := performPreCommitOperations(exe, settings, pr, opts)
	if err != nil {
		return err
	}

	if prID != "" {
		// If the PR exists, update it by pushing to the remote
		return update(exe, pr.branchID, prID)
	}
	// If it doesn't exist, return an error
	return errors.New("No PR found for branch " + pr.branchID)
}
