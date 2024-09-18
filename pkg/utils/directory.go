package utils

import (
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

// SetWorkDirToGitHubRoot checks whether the current working directory is in a GitHub repo.
// If false, an error is raised. If true, the working directory is set to the root of that repo.
func SetWorkDirToGitHubRoot(exe Executor) error {
	_, err := isInGitHubRepo(exe)
	if err != nil {
		return err
	}

	root, err := GetGitRootDirectory(exe)
	if err != nil {
		return err
	}

	err = exe.Chdir(root)
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

	files = strings.TrimSpace(files)

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
