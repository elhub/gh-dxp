package diff

import (
	"fmt"
	"strings"

	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/pkg/errors"
)

func Execute(exe utils.Executor, options *Options) error {
	// Get branchID
	currentBranch, errBranch := exe.Command("git", "branch", "--show-current")
	if errBranch != nil {
		return errBranch
	}
	branchId := strings.Trim(currentBranch, "\n")

	// Check if PR exists on branch
	prId, errCheck := CheckForExistingPR(exe, branchId)
	if errCheck != nil {
		return errCheck
	}

	if prId != "" {
		// If the PR exists, update it by pushing to the remote
		return update(exe, branchId, prId)
	} else {
		// If it doesn't exist, create a new PR
		return create(exe, options, branchId)
	}
}

func create(exe utils.Executor, options *Options, branchId string) error {
	// Push the current branch to git remote
	s := utils.StartSpinner("Pushing current branch to remote...", "Pushed working branch to remote.")
	currentBranch, err := exe.Command("git", "push", "--set-upstream", "origin", branchId)
	s.Stop()
	if err != nil {
		return err
	}
	log.Info(strings.Trim(currentBranch, "\n"))
	fmt.Println("Current Branch: ", currentBranch)

	// Fetch the default branch
	baseBranch := options.baseBranch
	if baseBranch == "" {
		s := utils.StartSpinner("Fetching repository default branch...", "Fetched repository default branch")
		stdOut, err := exe.GH("repo", "view", "--json", "defaultBranchRef", "--jq", ".defaultBranchRef.name")
		s.Stop()
		if err != nil {
			return errors.Wrap(err, "Failed to fetch default branch")
		}
		baseBranch = strings.Trim(stdOut.String(), "\n")
	}

	pr, err := createPR(exe, options, branchId, baseBranch)
	if err != nil {
		return err
	}

	s = utils.StartSpinner("Processing pull request...", "Pull request "+pr.Title+" created.")
	args := []string{"pr", "create", "--title", pr.Title, "--body", pr.Body, "--base", baseBranch}
	args = append(args, generatePRArgs(options)...)
	stdOut, err := exe.GH(args...)
	s.Stop()
	if err != nil {
		return errors.Wrap(err, "Failed to create pull request")
	}
	log.Info(strings.Trim(stdOut.String(), "\n"))

	return nil
}

func generatePRArgs(options *Options) []string {
	args := []string{}

	if len(options.Assignees) > 0 {
		args = append(args, "--assignee", strings.Join(options.Assignees, ","))
	}
	if len(options.Reviewers) > 0 {
		args = append(args, "--reviewer", strings.Join(options.Reviewers, ","))
	}

	return args
}

func update(exe utils.Executor, branchId string, prId string) error {
	// Push the current branch to the already existing git remote
	s := utils.StartSpinner("Updating Pull Request #"+prId+"...", "Pull Request #"+prId+" has been updated.")
	_, err := exe.Command("git", "push")
	s.Stop()
	if err != nil {
		return err
	}

	// Fetching this for info
	stdOut, err := exe.GH("pr", "list", "-H", branchId, "--json", "url", "--jq", ".[].url")
	if err != nil {
		return err
	}

	log.Info(strings.Trim(stdOut.String(), "\n") + "\n")

	return nil
}
