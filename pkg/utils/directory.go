package utils

import (
	"os"
	"regexp"
	"strings"
)

// GetGitRootDirectory returns the root directory of the current git repo.
func GetGitRootDirectory(exe Executor) (string, error) {
	// Locate the root directory of current git repo
	// Fails if not in a repo
	root, err := exe.Command("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}

	formattedRoot := strings.TrimSuffix(root, "\n")

	return formattedRoot, nil
}

// SetWorkDirToGitHubRoot checks whether the current working directory is in a GitHub repo. If false, an error is raised. If true, the working directory is set to the root of that repo.
func SetWorkDirToGitHubRoot(exe Executor) error {
	_, err := isInGitHubRepo(exe)
	if err != nil {
		return err
	}
	err = setWorkingDirectoryToGitRoot(exe)
	return err
}

// setWorkingDirectoryToGitRoot sets the working directory to be the git root directory.
func setWorkingDirectoryToGitRoot(exe Executor) error {
	root, err := GetGitRootDirectory(exe)
	if err != nil {
		return err
	}

	err = os.Chdir(root)
	return err
}

// NotAGitHubRepoError signifies that the current working directory is not a GitHub repo.
type NotAGitHubRepoError struct {
	Msg string
}

// Signifies that the current working directory is not a GitHub repo.
func (e *NotAGitHubRepoError) Error() string {
	return e.Msg
}

// ListFilesInDirectory returns a list of files in a given directory.
func ListFilesInDirectory(exe Executor, directory string) ([]string, error) {
	// List all files in a directory
	// Fails if directory does not exist

	files, err := exe.Command("ls", directory)
	if err != nil {
		return nil, err
	}

	return strings.Split(files, "\n"), nil
}

// isInGitHubRepo checks whether the current working directory is in a GitHub repo.
func isInGitHubRepo(exe Executor) (bool, error) {
	url, err := exe.Command("git", "remote", "get-url", "origin")

	if err != nil {
		return false, err
	}
	if !urlIsGitHubRepo(url) {
		return false, &NotAGitHubRepoError{Msg: "Current origin is not a GitHub repository"}
	}

	return true, nil
}

func urlIsGitHubRepo(url string) bool {
	return (strings.HasPrefix(url, "https://github.com") || strings.HasPrefix(url, "git@github.com"))
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
