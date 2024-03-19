package utils

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/caarlos0/log"
)

type Executor struct {
	ExecCmd func(name string, arg ...string) *exec.Cmd
}

func Exec() *Executor {
	return &Executor{
		ExecCmd: exec.Command,
	}
}

func (e *Executor) Run(name string, args ...string) (string, error) {
	log.Debug(fmt.Sprintf("Running '%s %s'", name, strings.Join(args, " ")))
	cmd := e.ExecCmd(name, args...)
	bytes, err := cmd.CombinedOutput()
	return string(bytes), err
}
