package utils

import (
	"strings"
)

func GetGitRootDirectory(exe Executor) (string, error) {
	// Locate the root directory of current git repo
	// Fails if not in a repo

	root, err := exe.Command("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", &NotAGitRepoError{Msg: "Not a git repo"}
	}

	return strings.TrimSuffix(root, "\n"), nil
}

// NotAGitRepoError signifies that the current working directory is not a git repo.
type NotAGitRepoError struct {
	Msg string
}

// Signifies that the current working directory is not a git repo.
func (e *NotAGitRepoError) Error() string {
	return e.Msg
}
