// Package pr contains the functions and types for the pull request command.
package pr

import (
	"strings"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/lint"
	"github.com/elhub/gh-dxp/pkg/logger"
	"github.com/elhub/gh-dxp/pkg/test"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/pkg/errors"
)

// CheckForExistingPR checks if a PR already exists for the current branch.
func CheckForExistingPR(exe utils.Executor, branchID string) (string, error) {
	stdOut, err := exe.GH("pr", "list", "-H", branchID, "--json", "number", "--jq", ".[].number")

	if err != nil {
		logger.Debug("Error: " + err.Error())
		return "", errors.New("Failed to find existing PR")
	}

	number := strings.Trim(stdOut, "\n")

	return number, nil
}

// GetPRTitle gets the title of the current PR.
func GetPRTitle(exe utils.Executor) (string, error) {
	return getPRField(exe, "title")
}

// GetPRBody gets the body of the current PR.
func GetPRBody(exe utils.Executor) (string, error) {
	return getPRField(exe, "body")
}

// getPRField gets a single field of information about a PR using the gh pr command
func getPRField(exe utils.Executor, field string) (string, error) {
	stdOut, err := exe.GH("pr", "view", "--json", field, "--jq", "."+field)

	if err != nil {
		return "", err
	}

	value := strings.Trim(stdOut, "\n")

	return value, nil
}

func handleUncommittedChanges(exe utils.Executor, options *Options) ([]string, error) {
	// Handle presence of untracked changes - ignore or abort
	untrackedChanges, err := utils.GetUntrackedChanges(exe)
	if err != nil {
		return []string{}, err
	}

	if len(untrackedChanges) > 0 && !options.TestRun {
		res, err := utils.AskToConfirm(formatUntrackedFileChangesQuestion(untrackedChanges))
		if err != nil {
			return []string{}, err
		}
		if !res {
			return []string{}, errors.New("User aborted workflow")
		}
	}

	// Handle presence of tracked changes - commit or abort
	trackedChanges, err := utils.GetTrackedChanges(exe)
	if err != nil {
		return []string{}, err
	}

	if len(trackedChanges) > 0 && !options.TestRun && options.CommitMessage == "" {
		res, err := utils.AskToConfirm(formatTrackedFileChangesQuestion(trackedChanges))
		if err != nil {
			return []string{}, err
		}

		if !res {
			return []string{}, errors.New("User aborted workflow")
		}
	}
	return trackedChanges, nil
}

func addAndCommitFiles(exe utils.Executor, options *Options) error {
	var commitMessage string
	var err error

	if options.CommitMessage != "" {
		commitMessage = options.CommitMessage
	} else {

		if !options.TestRun {
			commitMessage, err = utils.AskForString("Please enter a commit message: ", "")
			if err != nil {
				return err
			} else if len(commitMessage) == 0 {
				return errors.New("Empty commit message not allowed")
			}
		} else {
			commitMessage = "default commit message"
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

func performPreCommitOperations(exe utils.Executor, settings *config.Settings, pr PullRequest, options *Options) (PullRequest, error) {
	// Handle uncommitted changes
	filesToCommit, err := handleUncommittedChanges(exe, options)
	if err != nil {
		return pr, err
	}

	// Run lint
	if !options.NoLint {
		err = lint.Run(exe, settings, &lint.Options{})
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
