// Package prmerge contains the logic for merging pull requests.
package prmerge

import (
	"errors"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/caarlos0/log"
	pr "github.com/elhub/gh-dxp/pkg/prcreate"
	"github.com/elhub/gh-dxp/pkg/utils"
)

// Execute merges a pull request on the current branch.
func Execute(exe utils.Executor, options *Options) error {
	// Get branchID
	currentBranch, errBranch := exe.Command("git", "branch", "--show-current")
	if errBranch != nil {
		return errBranch
	}
	branchID := strings.Trim(currentBranch, "\n")

	// Get prID
	prID, errPR := pr.CheckForExistingPR(exe, branchID)
	if errPR != nil {
		return errPR
	}

	prTitle, errTitle := pr.GetPRTitle(exe)
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
