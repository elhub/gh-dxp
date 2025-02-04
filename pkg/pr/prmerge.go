package pr

import (
	"errors"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/elhub/gh-dxp/pkg/logger"
	"github.com/elhub/gh-dxp/pkg/utils"
)

// ExecuteMerge merges a pull request on the current branch.
func ExecuteMerge(exe utils.Executor, options *MergeOptions) error {
	// Get branchID
	currentBranch, errBranch := exe.Command("git", "branch", "--show-current")
	if errBranch != nil {
		return errBranch
	}
	branchID := strings.Trim(currentBranch, "\n")

	// Get prID
	prID, errPR := CheckForExistingPR(exe, branchID)
	if errPR != nil {
		return errPR
	}

	prTitle, errTitle := GetPRTitle(exe)
	if errTitle != nil {
		return errTitle
	}

	prBody, errBody := GetPRBody(exe)
	if errBody != nil {
		return errBody
	}

	logger.Info("Merging pull request #" + prID + "(" + prTitle + ")")
	// TODO: Add list of commits
	doMerge := false
	if options.AutoConfirm {
		doMerge = true
	} else {
		errMerge := survey.AskOne(&survey.Confirm{
			Message: "Merge these changes?",
		}, &doMerge, survey.WithValidator(survey.Required))
		if errMerge != nil {
			return errMerge
		}
	}

	if !doMerge { // Exit
		return nil
	}

	stdOut, err := exe.GH("pr", "merge", "--squash", "--delete-branch", "--subject", prTitle, "--body", prBody)
	logger.Info(stdOut)

	if err != nil {
		logger.Debug("Error: " + err.Error())
		return errors.New("Failed to merge pull request #" + prID)
	}

	logger.Info("Deleted local " + branchID + " and switched to branch main")
	logger.Info("Deleted remote branch " + branchID)

	return nil
}
