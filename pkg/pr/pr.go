// Package pr contains the functions and types for the pull request command.
package pr

import (
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

func handleUncommittedChanges(exe ghutil.Executor, options *Options) ([]string, error) {
	// Handle presence of untracked changes - ignore or abort
	untrackedChanges, err := ghutil.GetUntrackedChanges(exe)
	if err != nil {
		return []string{}, err
	}

	if len(untrackedChanges) > 0 && !options.TestRun {
		res, err := ghutil.AskToConfirm(formatUntrackedFileChangesQuestion(untrackedChanges))
		if err != nil {
			return []string{}, err
		}
		if !res {
			return []string{}, errors.New("User aborted workflow")
		}
	}

	// Handle presence of tracked changes - commit or abort
	trackedChanges, err := ghutil.GetTrackedChanges(exe)
	if err != nil {
		return []string{}, err
	}

	commits, err := exe.Command("git", "log", "--oneline", "origin/main..")
	if err != nil {
		return []string{}, err
	}

	if len(trackedChanges) == 0 && commits == "" {
		return []string{}, errors.New("No tracked changes found, skipping commit")
	}

	if !options.TestRun && options.CommitMessage == "" {
		res, err := ghutil.AskToConfirm(formatTrackedFileChangesQuestion(trackedChanges))
		if err != nil {
			return []string{}, err
		}

		if !res {
			return []string{}, errors.New("User aborted workflow")
		}
	}
	return trackedChanges, nil
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

func performPreCommitOperations(exe ghutil.Executor, settings *config.Settings, pr PullRequest, options *Options) (PullRequest, error) {
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
