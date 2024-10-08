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

	// Check if PR exists on branch
	prID, errCheck := CheckForExistingPR(exe, pr.branchID)
	if errCheck != nil {
		return errCheck
	}
	if prID == "" {
		// If PR does not exist, return an error
		return errors.New("No PR found for branch " + pr.branchID)
	}

	prOpts := &Options{
		TestRun:       options.TestRun,
		NoLint:        options.NoLint,
		NoUnit:        options.NoUnit,
		CommitMessage: options.CommitMessage,
	}

	pr, err := performPreCommitOperations(exe, settings, pr, prOpts)
	if err != nil {
		return err
	}

	return update(exe, pr.branchID, prID)
}
