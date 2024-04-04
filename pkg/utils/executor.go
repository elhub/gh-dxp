package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/caarlos0/log"
	"github.com/cli/go-gh/v2"
)

type Executor interface {
	Command(name string, args ...string) (string, error)
	GH(args ...string) (bytes.Buffer, bytes.Buffer, error)
}

type LinuxExecutorImpl struct {
	ExecCommand func(name string, args ...string) *exec.Cmd
}

func LinuxExecutor() *LinuxExecutorImpl {
	return &LinuxExecutorImpl{
		ExecCommand: exec.Command,
	}
}

// OS command
func (e *LinuxExecutorImpl) Command(name string, args ...string) (string, error) {
	log.Debug(fmt.Sprintf("Running '%s %s'", name, strings.Join(args, " ")))
	cmd := e.ExecCommand(name, args...)
	bytes, err := cmd.CombinedOutput()
	return string(bytes), err
}

// GH Command
func (e *LinuxExecutorImpl) GH(args ...string) (bytes.Buffer, bytes.Buffer, error) {
	return gh.Exec(args...)
}
