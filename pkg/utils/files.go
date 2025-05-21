package utils

import (
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

	// If main does not exist, we need to check it out in order to perform the comparison (relevant for new TC agents)
	if !contains(branchList, "main") {
		_, err := exe.Command("git", "fetch", "origin", "main")
		if err != nil {
			return []string{}, err
		}
	}

	var changedFiles []string

	if len(branchList) > 0 {
		changedFilesString, err := exe.Command("git", "diff", "--name-only", "main", "--relative")
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
	re := regexp.MustCompile(`^([ADMRT]|\s)([ADMRT]|\s)\s`) // This regex is intended to catch all tracked changes except for unmerged conflicts
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

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
