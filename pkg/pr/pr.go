// Package pr contains the functions and types for the pull request command.
package pr

import (
	"path/filepath"
	"strings"

	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/lint"
	"github.com/elhub/gh-dxp/pkg/test"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/pkg/errors"
)

// CheckForExistingPR checks if a PR already exists for the current branch.
func CheckForExistingPR(exe utils.Executor, branchID string) (string, error) {
	stdOut, err := exe.GH("pr", "list", "-H", branchID, "--json", "number", "--jq", ".[].number")

	if err != nil {
		log.Debug("Error: " + err.Error())
		return "", errors.New("Failed to find existing PR")
	}

	number := strings.Trim(stdOut.String(), "\n")

	return number, nil
}

// GetPRTitle gets the title of the current PR.
func GetPRTitle(exe utils.Executor) (string, error) {
	stdOut, err := exe.GH("pr", "view", "--json", "title", "--jq", ".title")

	if err != nil {
		return "", errors.New("Error getting PR title")
	}

	title := strings.Trim(stdOut.String(), "\n")

	return title, nil
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

func addAndCommitFiles(exe utils.Executor, files []string, options *Options) error {
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
	// Get git root directory and add to files to get fully qualified paths
	root, err := utils.GetGitRootDirectory(exe)
	if err != nil {
		return err
	}

	var fullPaths []string
	for _, filePath := range files {
		fullPaths = append(fullPaths, filepath.Join(root, filePath))
	}

	addCommandArgs := append([]string{"add"}, fullPaths...)

	_, err = exe.Command("git", addCommandArgs...)
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
		err = addAndCommitFiles(exe, filesToCommit, options)
		if err != nil {
			return pr, err
		}
	}

	return pr, nil
}
