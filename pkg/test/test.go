package test

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/utils"
)

// FileExists checks to see whether a file exists in the file system. Exported to allow mocking during tests.
var FileExists = utils.FileExists

// RunTest runs a workflow to automatically determine relevant tests in the current repo and run them.
func RunTest(exe utils.Executor) error {
	cmd, err := resolveTestCommand(exe)
	if err != nil {
		var ntcErr *NoTestCommandError
		if errors.As(err, &ntcErr) {
			log.Warn(err.Error())
			return nil
		}
		return err
	}

	ctx := context.Background()
	err = exe.CommandContext(ctx, cmd, "test")
	if err != nil {
		return err
	}

	return nil
}

func resolveTestCommand(exe utils.Executor) (string, error) {
	root, err := getGitRootDirectory(exe)
	if err != nil {
		return "", err
	}

	if makeTestInGitRoot(root) {
		return "make", nil
	}

	if gradleTestInGitRoot(root) {
		return "./gradlew", nil
	}
	if mavenTestInGitRoot(root) {
		return "mvn", nil
	}
	if npmTestInGitRoot(root) {
		return "npm", nil
	}
	return "", &NoTestCommandError{Msg: "No test command found"}
}

func getGitRootDirectory(exe utils.Executor) (string, error) {
	// Locate the root directory of current git repo
	// Fails if not in a repo

	root, err := exe.Command("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", &NotAGitRepoError{Msg: "Not a git repo"}
	}

	return strings.TrimSuffix(root, "\n"), nil
}

func gradleTestInGitRoot(root string) bool {
	return FileExists(fmt.Sprintf("%s/gradlew", root))
}

func makeTestInGitRoot(root string) bool {
	return FileExists(fmt.Sprintf("%s/Makefile", root))
}

func mavenTestInGitRoot(root string) bool {
	return FileExists(fmt.Sprintf("%s/pom.xml", root))
}

func npmTestInGitRoot(root string) bool {
	return FileExists(fmt.Sprintf("%s/package.json", root))
}

// NoTestCommandError signifies that no valid test command was found in the current git repo.
type NoTestCommandError struct {
	Msg string
}

// NotAGitRepoError signifies that the current working directory is not a git repo.
type NotAGitRepoError struct {
	Msg string
}

// Signifies that the current working directory is not a git repo.
func (e *NotAGitRepoError) Error() string {
	return e.Msg
}

// Signifies that no valid test command was found in the current git repo.
func (e *NoTestCommandError) Error() string {
	return e.Msg
}
