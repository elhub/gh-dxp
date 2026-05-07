// Package pr contains the functions and types for the pull request command.
package pr

import (
	"fmt"
	"strings"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/ghutil"
	"github.com/elhub/gh-dxp/pkg/lint"
	"github.com/elhub/gh-dxp/pkg/logger"
	"github.com/elhub/gh-dxp/pkg/renovate"
	"github.com/elhub/gh-dxp/pkg/test"
	"github.com/pkg/errors"
)

// CheckForExistingPR checks if a PR already exists for the current branch.
func CheckForExistingPR(exe ghutil.Executor, branchID string) (string, error) {
	stdOut, err := exe.GH("pr", "list", "-H", branchID, "--json", "number", "--jq", ".[].number")

	if err != nil {
		logger.Debug("Error: " + err.Error())
		return "", errors.New("Failed to find existing PR")
	}

	number := strings.Trim(stdOut, "\n")

	return number, nil
}

// GetPRTitle gets the title of the current PR.
func GetPRTitle(exe ghutil.Executor) (string, error) {
	return getPRField(exe, "title")
}

// GetPRBody gets the body of the current PR.
func GetPRBody(exe ghutil.Executor) (string, error) {
	return getPRField(exe, "body")
}

// getPRField gets a single field of information about a PR using the gh pr command.
func getPRField(exe ghutil.Executor, field string) (string, error) {
	stdOut, err := exe.GH("pr", "view", "--json", field, "--jq", "."+field)

	if err != nil {
		return "", err
	}

	value := strings.Trim(stdOut, "\n")

	return value, nil
}

// ValidateLocalChanges checks for untracked, uncommitted, and committed changes in the local repository.
// It returns a list of uncommitted tracked changes that should be included in the PR.
func validateLocalChanges(exe ghutil.Executor, options *Options) ([]string, error) {
	if err := handleUntrackedChanges(exe, options); err != nil {
		return []string{}, err
	}

	uncommittedChanges, err := handleUncommittedChanges(exe, options)
	if err != nil {
		return []string{}, err
	}

	committedChanges, err := handleCommittedChanges(exe)
	if err != nil {
		return []string{}, err
	}

	if len(uncommittedChanges) == 0 && len(committedChanges) == 0 {
		return []string{}, errors.New("No changes found, aborting PR operation")
	}

	return uncommittedChanges, nil
}

func handleUntrackedChanges(exe ghutil.Executor, options *Options) error {
	untrackedChanges, err := ghutil.GetUntrackedChanges(exe)
	if err != nil {
		return err
	}

	// Skip if no untracked changes or in test mode
	if len(untrackedChanges) == 0 || options.TestRun {
		return nil
	}

	confirmed, err := ghutil.AskToConfirm(formatUntrackedFileChangesQuestion(untrackedChanges))
	if err != nil {
		return err
	}

	if !confirmed {
		return errors.New("User aborted workflow")
	}

	return nil
}

func handleCommittedChanges(exe ghutil.Executor) ([]string, error) {
	commits, err := exe.Command("git", "log", "--oneline", "origin/main..")
	if err != nil {
		return []string{}, err
	}
	commitsList := []string{}
	if commits != "" {
		commits = strings.TrimSpace(commits)
		commitsList = strings.Split(commits, "\n")
		logger.Info("Using already committed changes for PR:")
		for _, commit := range commitsList {
			if len(commit) > 0 {
				logger.Info(fmt.Sprintf("\t- %s", commit))
			}
		}
		logger.Info("")
	} else {
		logger.Debug("No committed changes found")
	}
	return commitsList, nil
}

func handleUncommittedChanges(exe ghutil.Executor, options *Options) ([]string, error) {
	uncommittedTrackedChanges, err := ghutil.GetTrackedChanges(exe)
	if err != nil {
		return []string{}, err
	}

	if len(uncommittedTrackedChanges) == 0 {
		logger.Debug("No uncommitted tracked changes found")
		return []string{}, nil
	}

	// Skip confirmation if test run or commit message already provided
	if options.TestRun || options.CommitMessage != "" {
		return uncommittedTrackedChanges, nil
	}

	if len(uncommittedTrackedChanges) > 0 {
		confirmed, err := ghutil.AskToConfirm(formatTrackedFileChangesQuestion(uncommittedTrackedChanges))
		if err != nil {
			return []string{}, err
		}
		if !confirmed {
			return []string{}, errors.New("User aborted workflow")
		}
	}

	return uncommittedTrackedChanges, nil
}

func addAndCommitFiles(exe ghutil.Executor, options *Options) error {
	var commitMessage string
	var err error

	switch {
	case options.CommitMessage != "":
		commitMessage = options.CommitMessage
	case options.TestRun:
		commitMessage = "default commit message"
	default:
		commitMessage, err = ghutil.AskForString("Please enter a commit message: ", "")
		if err != nil {
			return err
		}
		if len(commitMessage) == 0 {
			return errors.New("Empty commit message not allowed")
		}
	}

	_, err = exe.Command("git", "add", "-u")
	if err != nil {
		return err
	}

	// Commit files
	_, err = exe.Command("git", "commit", "-m", commitMessage)
	if err != nil {
		return err
	}

	return nil
}

func performPreCreateOperations(exe ghutil.Executor, settings *config.Settings, pr PullRequest, options *Options) (PullRequest, error) {
	// Handle uncommitted changes
	filesToCommit, err := validateLocalChanges(exe, options)
	if err != nil {
		return pr, err
	}

	// Run lint
	if !options.NoLint {
		err = lint.Run(exe, settings, &lint.Options{})
		if err != nil {
			return pr, err
		}

		// Run renovate config validation
		err = renovate.Run(exe, settings, &renovate.Options{})
		if err != nil {
			return pr, err
		}
		pr.isLinted = true
	}

	// Run tests
	if !options.NoUnit {
		pr.isTested, err = test.RunTest(exe)
		if err != nil {
			return pr, err
		}
	}

	if len(filesToCommit) > 0 {
		err = addAndCommitFiles(exe, options)
		if err != nil {
			return pr, err
		}
	}

	return pr, nil
}
