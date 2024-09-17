// Package utils provides common utilities for the gh-dxp extension.
package utils

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/caarlos0/log"
	"github.com/cli/go-gh/v2"
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
	log.Debug(fmt.Sprintf("Running '%s %s'", name, strings.Join(args, " ")))
	cmd := e.ExecCommand(name, args...)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(bytes))
	}
	return string(bytes), err
}

// CommandContext runs an OS command with a context and returns an error.
// The output is printed to stdout and stderr.
func (e *LinuxExecutorImpl) CommandContext(ctx context.Context, name string, args ...string) error {
	log.Debug(fmt.Sprintf("Running with context '%s %s'", name, strings.Join(args, " ")))
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

// GH runs a GitHub CLI command and returns its output.
func (e *LinuxExecutorImpl) GH(args ...string) (bytes.Buffer, error) {
	log.Debug(fmt.Sprintf("Running gh '%s'", strings.Join(args, " ")))
	stdOut, stdErr, err := gh.Exec(args...)
	if err != nil {
		log.Error(stdErr.String())
		log.Debug(fmt.Sprintf("Error running GH command: %s", err.Error()))
		return stdErr, err
	}
	return stdOut, err
}

// Chdir changes the current working directory.
func (e *LinuxExecutorImpl) Chdir(dir string) error {
	return os.Chdir(dir)
}
