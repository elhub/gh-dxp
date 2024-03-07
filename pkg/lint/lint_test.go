package lint_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/elhub/gh-devxp/pkg/config"
	"github.com/elhub/gh-devxp/pkg/lint"
	"github.com/elhub/gh-devxp/pkg/utils"
)

// Test golint.go
func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestGoLintErrors", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestGoLintErrors(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintf(os.Stdout, "pkg/test/main.go:17:15: error returned from external package is unwrapped: sig: func gopkg.in/yaml.v3.Unmarshal(in []byte, out interface{}) (err error) (wrapcheck)\n")
	fmt.Fprintf(os.Stdout, "                return nil, err\n")
	fmt.Fprintf(os.Stdout, "                            ^\n")
	fmt.Fprintf(os.Stdout, "pkg/lint/main.go:25: unnecessary trailing newline (whitespace)\n")
	fmt.Fprintf(os.Stdout, "                             \n")
	fmt.Fprintf(os.Stdout, "                } else {     \n")
	os.Exit(0)
}

func TestGoLint(t *testing.T) {
	utils.ExecCmd = mockExecCommand
	defer func() { utils.ExecCmd = exec.Command }() // Restore after test

	// Call GoLint.Exec
	outputs, err := lint.GoLint{}.Exec()

	if err != nil {
		t.Errorf("GoLint threw an error when none was expected")
	}

	if len(outputs) != 2 {
		t.Errorf("Expected 2 issues in GoLint.Exec, got '%d'", len(outputs))
	}

	if outputs[0].Linter != "golint" {
		t.Errorf("Expected first issue linter to be 'golint', got '%s'", outputs[0].Linter)
	}

	if outputs[0].Path != "pkg/test/main.go" {
		t.Errorf("Expected first issue path to be 'pkg/test/main.go', got '%s'", outputs[0].Path)
	}

}

// Test main.go
type TestMockLint struct{}

func (TestMockLint) Exec() ([]lint.LintOutput, error) {
	return []lint.LintOutput{
		{
			Linter:      "golint",
			Path:        "mock/path",
			Line:        1,
			Character:   1,
			Code:        "mock/code",
			Description: "mock/description",
		},
	}, nil
}

func TestRun(t *testing.T) {
	lint.Linters = map[string]lint.Linter{
		"golint": TestMockLint{},
	}

	// Create a mock context and settings
	ctx := context.Background()
	settings := &config.Settings{
		Lint: config.LintSettings{
			Linters: []config.LinterSettings{
				{
					Name: "mocklint",
				},
			},
		},
	}

	// Call Run with the mock context and settings
	err := lint.Run(ctx, settings)

	// Check if Run returns an error
	if err != nil {
		t.Errorf("Expected Run to not return an error, got %v", err)
	}

}
