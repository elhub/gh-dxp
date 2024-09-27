package pr

import (
	"errors"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/caarlos0/log"
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

	log.Info("Merging pull request #" + prID + "(" + prTitle + ")")
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

	stdOut, err := exe.GH("pr", "merge", "--squash", "--delete-branch")
	log.Info(stdOut.String())

	if err != nil {
		log.Debug("Error: " + err.Error())
		return errors.New("Failed to merge pull request #" + prID)
	}

	log.Info("Deleted local " + branchID + " and switched to branch main")
	log.Info("Deleted remote branch " + branchID)

	return nil
}
