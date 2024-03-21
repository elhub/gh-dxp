package branch

import (
	"os/exec"
	"strings"

	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/pkg/errors"
)

func CheckoutBranch(exe utils.Executor, branchID string) error {
	// Does the branch exist?
	branchExists, existsErr := branchExists(exe, branchID)
	if existsErr != nil {
		return existsErr
	}

	if branchExists {
		log.Debugf("Branch '%s' already exists, checking out to it", branchID)
		out, err1 := exe.Command("git", "checkout", branchID)
		if err1 != nil {
			return errors.Wrap(err1, "Failed to checkout branch")
		}
		log.Info(strings.Trim(out, "\n"))
	} else {
		log.Debugf("Creating branch '%s' and checking out to it", branchID)
		out, err2 := exe.Command("git", "checkout", "-b", branchID)
		if err2 != nil {
			return errors.Wrap(err2, "Failed to create branch")
		}
		log.Info(strings.Trim(out, "\n"))
	}

	return nil
}

func branchExists(exe utils.Executor, branchID string) (bool, error) {
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
