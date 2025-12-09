// Package utils provides command execution utilities for Linux systems.
//
//revive:disable
package utils

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

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

	// Set process group so we can kill all child processes
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// Start the command instead of Run() so we can handle context cancellation
	if err := cmd.Start(); err != nil {
		return err
	}

	// Wait for either the command to finish or context to be cancelled
	errChan := make(chan error, 1)
	go func() {
		errChan <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		// Context was cancelled, kill the entire process group
		if cmd.Process != nil {
			// Kill the process group (negative PID kills the group)
			if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
				logger.Error(fmt.Sprintf("Failed to kill process group %d: %v", cmd.Process.Pid, err))
			}
		}
		// Wait for the command to finish with a timeout
		select {
		case <-errChan:
			// Process finished
		case <-time.After(5 * time.Second):
			logger.Error(fmt.Sprintf("Process %d did not terminate within 5 seconds after SIGKILL", cmd.Process.Pid))
		}
		return ctx.Err()
	case err := <-errChan:
		return err
	}
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
