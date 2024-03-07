package utils

import (
	"os"
	"os/exec"
	"testing"
)

// mockExecCommand is a helper function that creates a mock of exec.Command
func mockExecCommand(command string, args ...string) *exec.Cmd {
	// This test binary always exits with 0 and prints whatever is passed to -echo
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// TestHelperProcess is a helper function for mockExecCommand
func TestHelperProcess(t *testing.T) {
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

func TestExec(t *testing.T) {
	// Replace exec.Command with our mock
	ExecCmd = mockExecCommand
	defer func() { ExecCmd = exec.Command }() // Restore after test

	// Call Exec
	_, err := Exec("echo", "hello")

	// Check if Exec returns an error
	if err != nil {
		t.Errorf("Expected Exec to not return an error, got %v", err)
	}
}
