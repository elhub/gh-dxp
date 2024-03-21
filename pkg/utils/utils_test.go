package utils_test

import (
	"os/exec"
	"testing"

	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestLinuxExecutor_Command(t *testing.T) {
	t.Run("should return command output", func(t *testing.T) {
		executor := &utils.LinuxExecutorImpl{
			ExecCommand: func(name string, args ...string) *exec.Cmd {
				return exec.Command("echo", "hello")
			},
		}

		// Call the method under test
		output, err := executor.Command("echo", "hello")

		// Assert that the expectations were met
		assert.NoError(t, err)
		assert.Equal(t, "hello\n", output)
	})

	t.Run("should return error if command fails", func(t *testing.T) {
		executor := &utils.LinuxExecutorImpl{
			ExecCommand: func(name string, args ...string) *exec.Cmd {
				return exec.Command("ls", "/nonexistent")
			},
		}

		// Call the method under test
		output, err := executor.Command("ls", "/nonexistent")

		// Assert that the expectations were met
		assert.Error(t, err)
		assert.Contains(t, output, "No such file or directory")
	})
}
