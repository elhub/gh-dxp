package utils

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/caarlos0/log"
)

func Exec(name string, args ...string) (string, error) {
	log.Debug(fmt.Sprintf("Running '%s %s'", name, strings.Join(args, " ")))
	cmd := exec.Command(name, args...)
	bytes, err := cmd.CombinedOutput()
	return string(bytes), err
}
