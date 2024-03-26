package utils

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/caarlos0/log"
)

type Executor interface {
	Command(name string, args ...string) (string, error)
}

type LinuxExecutorImpl struct {
	ExecCommand func(name string, args ...string) *exec.Cmd
}

func LinuxExecutor() *LinuxExecutorImpl {
	return &LinuxExecutorImpl{
		ExecCommand: exec.Command,
	}
}

func (e *LinuxExecutorImpl) Command(name string, args ...string) (string, error) {
	log.Debug(fmt.Sprintf("Running '%s %s'", name, strings.Join(args, " ")))
	cmd := e.ExecCommand(name, args...)
	bytes, err := cmd.CombinedOutput()
	return string(bytes), err
}
