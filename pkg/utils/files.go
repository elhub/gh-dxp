// Package utils provides common utilities for the gh-dxp extension.
package utils

import (
	"github.com/elhub/gh-dxp/pkg/logger"
	"os"
	"regexp"
	"strings"
)

// FileExists returns true if a given file exists, and false if it doesn't.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetChangedFiles returns a list of changed files in the current repo
func GetChangedFiles(exe Executor) ([]string, error) {
	branchString, err := exe.Command("git", "branch")
	if err != nil {
		return []string{}, err
	}

	branchList := ConvertTerminalOutputIntoList(branchString)

	var changedFiles []string

	if len(branchList) > 0 {

		// Pull latest changes from main branch
		logger.Info("Fetching latest changes from the main branch...")
		_, pullErr := exe.Command("git", "fetch", "origin", "main")
		if pullErr != nil {
			return nil, err
		}

		// Locally update origin/HEAD symbolic reference to point at the default branch
		_, err := exe.Command("git", "remote", "set-head", "origin", "--auto")
		if err != nil {
			return nil, err
		}

		// Fetch the default branch (should be the main branch or a temporary working branch that’s intended to be merged back into main)
		headRef, err := exe.Command("git", "symbolic-ref", "--short", "refs/remotes/origin/HEAD")
		if err != nil {
			return nil, err
		}

		headRef = strings.TrimSpace(headRef)
		logger.Info("Checking for changes relative to the default branch: " + headRef)
		changedFilesString, err := exe.Command("git", "diff", "--name-only", headRef, "--relative")
		changedFiles = ConvertTerminalOutputIntoList(changedFilesString)
		if err != nil {
			return []string{}, err
		}
	} else {
		changedFiles, err = GetTrackedChanges(exe)
		if err != nil {
			return []string{}, err
		}
	}
	return changedFiles, nil
}

// CheckFilesUpdated checks if any of the specified patterns match the changed files.
func CheckFilesUpdated(changedFiles []string, patterns []string) bool {
	for _, file := range changedFiles {
		for _, pattern := range patterns {
			matched, err := regexp.MatchString(pattern, file)
			if err != nil {
				continue
			}
			if matched {
				return true
			}
		}
	}
	return false
}

// GetUntrackedChanges returns a list of file names for unchanged files in the current repo
func GetUntrackedChanges(exe Executor) ([]string, error) {
	re := regexp.MustCompile(`^\?\?`)

	return getChanges(exe, re)
}

// GetTrackedChanges returns a list of file names for changed files in the current repo
func GetTrackedChanges(exe Executor) ([]string, error) {
	// This regex is intended to catch all tracked changes except for unmerged conflicts
	// We need to check for strings like ' M example` and `M  example` to catch both staged and unstaged changes.
	re := regexp.MustCompile(`^([ADMRT]|\s)([ADMRT]|\s)\s`)
	return getChanges(exe, re)
}

// Checks the current repo state for any changes matching a given regex 're'
func getChanges(exe Executor, re *regexp.Regexp) ([]string, error) {
	changeString, err := exe.Command("git", "status", "--porcelain")
	if err != nil {
		return []string{}, err
	}

	changes := strings.Split(changeString, "\n")
	matchedChanges := filter(changes, re.MatchString)

	// Remove the regex matched part of the string, leaving only the file name
	for i, s := range matchedChanges {
		matchedChanges[i] = re.ReplaceAllString(s, "")
	}

	// Split string on '->' and capture the last element, in order to catch changes of type 'oldfilename.txt -> newfilename.txt'
	for i, s := range matchedChanges {
		macgyver := strings.Split(s, "->")
		matchedChanges[i] = strings.TrimSpace(macgyver[len(macgyver)-1])
	}
	return matchedChanges, nil
}

func filter(list []string, test func(string) bool) []string {
	ret := []string{}
	for _, s := range list {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return ret
}
