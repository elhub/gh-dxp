package merge

import (
	"errors"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/diff"
	"github.com/elhub/gh-dxp/pkg/utils"
)

func Execute(exe utils.Executor, options *Options) error {
	// Get branchID
	currentBranch, errBranch := exe.Command("git", "branch", "--show-current")
	if errBranch != nil {
		return errBranch
	}
	branchId := strings.Trim(currentBranch, "\n")

	// Get prID
	prId, errPr := diff.CheckForExistingPR(exe, branchId)
	if errPr != nil {
		return errPr
	}

	prTitle, errTitle := diff.GetPRTitle(exe)
	if errTitle != nil {
		return errTitle
	}

	log.Info("Merging pull request #" + prId + "(" + prTitle + ")")
	// TODO: Add list of commits
	doMerge := false
	if options.AutoConfirm {
		doMerge = true
	} else {
		survey.AskOne(&survey.Confirm{
			Message: "Merge these changes?",
		}, &doMerge, survey.WithValidator(survey.Required))
	}

	if !doMerge { // Exit
		return nil
	}

	stdOut, err := exe.GH("pr", "merge", "--squash", "--auto", "--delete-branch")
	log.Info(stdOut.String())

	if err != nil {
		log.Debug("Error: " + err.Error())
		return errors.New("Failed to merge pull request #" + prId)
	}

	log.Info("Deleted local " + branchId + " and switched to branch main")
	log.Info("Deleted remote branch " + branchId)

	return nil
}
