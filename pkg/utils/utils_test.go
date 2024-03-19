package utils_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/elhub/gh-devxp/pkg/utils"
)

// mockExecCommand is a helper function that creates a mock of exec.Command.
func mockExecCommand(command string, args ...string) *exec.Cmd {
	// This test binary always exits with 0 and prints whatever is passed to -echo.
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// TestHelperProcess is a helper function for mockExecCommand.
func TestHelperProcess(_ *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// Echo whatever comes after -echo
	for _, arg := range os.Args {
		if arg == "-echo" {
			println(arg)
			os.Exit(0)
		}
	}
}

func mockExec() *utils.Executor {
	return &utils.Executor{
		ExecCmd: mockExecCommand,
	}
}

func TestExec(t *testing.T) {
	_, err := mockExec().Run("echo", "hello")
	if err != nil {
		t.Errorf("Expected Exec to not return an error, got %v", err)
	}
}
