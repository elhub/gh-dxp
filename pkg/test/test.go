package test

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/utils"
)

// FileExists checks to see whether a file exists in the file system.
var FileExists = utils.FileExists //nolint: gochecknoglobals // Exported to allow mocking during tests.

// RunTest runs a workflow to automatically determine relevant tests in the current repo and run them.
func RunTest(exe utils.Executor) error {
	cmd, args, err := resolveTestCommand(exe)
	if err != nil {
		var ntcErr *NoTestCommandError
		if errors.As(err, &ntcErr) {
			log.Warn(err.Error())
			return nil
		}
		return err
	}

	ctx := context.Background()
	err = exe.CommandContext(ctx, cmd, args...)
	if err != nil {
		return err
	}

	return nil
}

func resolveTestCommand(exe utils.Executor) (string, []string, error) {
	root, err := utils.GetGitRootDirectory(exe)
	if err != nil {
		return "", nil, err
	}

	if makeTestInGitRoot(root) {
		return "make", []string{"check"}, nil
	}

	if gradleTestInGitRoot(root) {
		return "./gradlew", []string{"test"}, nil
	}
	if mavenTestInGitRoot(root) {
		return "mvn", []string{"test"}, nil
	}
	if npmTestInGitRoot(root) {
		return "npm", []string{"test"}, nil
	}
	return "", []string{}, &NoTestCommandError{Msg: "No test command could be automatically detected. If you want to automatically run tests as part of pr creation, please have a look at the documentation."}
}

func gradleTestInGitRoot(root string) bool {
	return FileExists(filepath.Join(root, ".gradlew"))
}

func makeTestInGitRoot(root string) bool {
	return FileExists(filepath.Join(root, "Makefile"))
}

func mavenTestInGitRoot(root string) bool {
	return FileExists(filepath.Join(root, "pom.xml"))
}

func npmTestInGitRoot(root string) bool {
	return FileExists(filepath.Join(root, "package.json"))
}

// NoTestCommandError signifies that no valid test command was found in the current git repo.
type NoTestCommandError struct {
	Msg string
}

// Signifies that no valid test command was found in the current git repo.
func (e *NoTestCommandError) Error() string {
	return e.Msg
}
