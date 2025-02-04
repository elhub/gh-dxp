// Package branch provides functions to work with git branches.
package branch

import (
	"os/exec"
	"strings"

	"github.com/elhub/gh-dxp/pkg/logger"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/pkg/errors"
)

// CheckoutBranch checks out to the branch with the given ID.
func CheckoutBranch(exe utils.Executor, branchID string) error {
	// Does the branch exist?
	branchExists, existsErr := Exists(exe, branchID)
	if existsErr != nil {
		return existsErr
	}

	if branchExists {
		logger.Debugf("Branch '%s' already exists, checking out to it", branchID)
		out, err1 := exe.Command("git", "checkout", branchID)
		if err1 != nil {
			return errors.Wrap(err1, "Failed to checkout branch")
		}
		logger.Info(strings.Trim(out, "\n"))
	} else {
		logger.Debugf("Creating branch '%s' and checking out to it", branchID)
		out, err2 := exe.Command("git", "checkout", "-b", branchID)
		if err2 != nil {
			return errors.Wrap(err2, "Failed to create branch")
		}
		logger.Debug(strings.Trim(out, "\n"))
	}

	return nil
}

// Exists checks whether a specified branch exists.
func Exists(exe utils.Executor, branchID string) (bool, error) {
	_, err := exe.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branchID)
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if exitErr.ExitCode() == 0 {
				return true, nil
			}
			return false, nil
		}

		return false, err
	}
	return true, nil
}

// GetCommitMessages returns the commit messages between the main branch and the branch with the given ID.
func GetCommitMessages(exe utils.Executor, mainID string, branchID string) (string, error) {
	return exe.Command("git", "log", mainID+".."+branchID, "--oneline", "--pretty=format:%s")
}
