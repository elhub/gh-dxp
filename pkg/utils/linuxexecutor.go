// Package utils provides common utilities for the gh-dxp extension.
package utils

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cli/go-gh/v2"
	"github.com/elhub/gh-dxp/pkg/logger"
)

// LinuxExecutorImpl is the type of the Executor interface for Linux systems.
type LinuxExecutorImpl struct {
	ExecCommand func(name string, args ...string) *exec.Cmd
}

// LinuxExecutor returns a new LinuxExecutorImpl.
func LinuxExecutor() *LinuxExecutorImpl {
	return &LinuxExecutorImpl{
		ExecCommand: exec.Command,
	}
}

// Command runs an OS command and returns its output.
func (e *LinuxExecutorImpl) Command(name string, args ...string) (string, error) {
	logger.Debug(fmt.Sprintf("Running '%s %s'", name, strings.Join(args, " ")))
	cmd := e.ExecCommand(name, args...)
	bytes, err := cmd.CombinedOutput()

	outputString := string(bytes)

	if err != nil && outputString != "" {
		logger.Error(outputString)
	}
	return string(bytes), err
}

// CommandContext runs an OS command with a context and returns an error.
// The output is printed to stdout and stderr.
func (e *LinuxExecutorImpl) CommandContext(ctx context.Context, name string, args ...string) error {
	logger.Debug(fmt.Sprintf("Running with context '%s %s'", name, strings.Join(args, " ")))
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

// GH runs a GitHub CLI command and returns its output.
func (e *LinuxExecutorImpl) GH(args ...string) (string, error) {
	logger.Debug(fmt.Sprintf("Running gh '%s'", strings.Join(args, " ")))
	stdOut, stdErr, err := gh.Exec(args...)
	if err != nil {
		logger.Error(stdErr.String())
		logger.Debug(fmt.Sprintf("Error running GH command: %s", err.Error()))
		return stdErr.String(), err
	}
	return stdOut.String(), err
}

// Chdir changes the current working directory.
func (e *LinuxExecutorImpl) Chdir(dir string) error {
	return os.Chdir(dir)
}
