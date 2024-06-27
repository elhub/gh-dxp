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
		return "", &NotAGitRepoError{Msg: "Not a git repo"}
	}

	formattedRoot := strings.TrimSuffix(root, "\n") + "/"

	return formattedRoot, nil
}

// NotAGitRepoError signifies that the current working directory is not a git repo.
type NotAGitRepoError struct {
	Msg string
}

// Signifies that the current working directory is not a git repo.
func (e *NotAGitRepoError) Error() string {
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
