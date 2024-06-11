package test

import (
	"context"
	"fmt"
	"strings"

	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/utils"
)

var FileExists = utils.FileExists

func RunTest(exe utils.Executor) error {

	cmd, err := resolveTestCommand(exe)
	if err != nil {
		switch e := err.(type) {
		case *NoTestCommandError:
			log.Warn(e.Msg)
			return nil
		default:
			return err
		}

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
	//Locate the root directory of current git repo
	//Fails if not in a repo

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

type NoTestCommandError struct {
	Msg string
}

type NotAGitRepoError struct {
	Msg string
}

func (e *NotAGitRepoError) Error() string {
	return e.Msg
}

func (e *NoTestCommandError) Error() string {
	return e.Msg
}
