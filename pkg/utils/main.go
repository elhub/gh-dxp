package utils

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/caarlos0/log"
)

var ExecCmd = exec.Command

func Exec(name string, args ...string) (string, error) {
	log.Debug(fmt.Sprintf("Running '%s %s'", name, strings.Join(args, " ")))
	cmd := ExecCmd(name, args...)
	bytes, err := cmd.CombinedOutput()
	return string(bytes), err
}
