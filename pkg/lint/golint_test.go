package lint_test

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/elhub/gh-devxp/pkg/lint"
	"github.com/elhub/gh-devxp/pkg/utils"
)

func TestGoLintParser(t *testing.T) {
	tests := []struct {
		name      string
		inputLine string
		expected  lint.LintOutput
		expectErr bool
	}{
		{
			name:      "valid golint line with line number and column",
			inputLine: "pkg/config/config_test.go:32:5: shadow declaration of \"err\" shadows declaration at line 10 (govet)",
			expectErr: false,
			expected: lint.LintOutput{
				Linter:      "golint",
				Path:        "pkg/config/config_test.go",
				Line:        32,
				Column:      5,
				Description: "shadow declaration of \"err\" shadows declaration at line 10",
				Severity:    "error",
				Source:      "govet",
			},
		},
		{
			name:      "valid golint line with line number and no column",
			inputLine: "config_test.go:66: unnecessary trailing newline (whitespace)",
			expectErr: false,
			expected: lint.LintOutput{
				Linter:      "golint",
				Path:        "config_test.go",
				Line:        66,
				Column:      0,
				Description: "unnecessary trailing newline",
				Severity:    "error",
				Source:      "whitespace",
			},
		},
		{
			name:      "valid golint line with colon in description",
			inputLine: "pkg/config/config_test.go:789:91: shadow: declaration of \"err\" shadows declaration at line 10 (govet)",
			expectErr: false,
			expected: lint.LintOutput{
				Linter:      "golint",
				Path:        "pkg/config/config_test.go",
				Line:        789,
				Column:      91,
				Description: "shadow: declaration of \"err\" shadows declaration at line 10",
				Severity:    "error",
				Source:      "govet",
			},
		},
		{
			name:      "invalid delimiters in golint format",
			inputLine: "[error] test:1,2 no new line character at the end of file (new-line-at-end-of-file)",
			expectErr: true,
			expected:  lint.LintOutput{},
		},
		{
			name:      "invalid line fomrat in golint",
			inputLine: "test.yaml:t:8: [error] no new line character at the end of file (new-line-at-end-of-file)",
			expectErr: true,
			expected:  lint.LintOutput{},
		},
		{
			name:      "invalid column format in golint",
			inputLine: "test.yaml:1:t: [error] no new line character at the end of file (new-line-at-end-of-file)",
			expectErr: true,
			expected:  lint.LintOutput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := lint.GoLintParser(tt.inputLine)
			if (err != nil) != tt.expectErr {
				t.Errorf("GoLintParser() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !reflect.DeepEqual(output, tt.expected) {
				t.Errorf("golintParser() output = %v, expected %v", output, tt.expected)
			}
		})
	}

}

// Test golint.go
func mockGoExec(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestGoLintErrors", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestGoLintErrors(_ *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintf(os.Stdout, "pkg/lint/golint.go:17:2: use of `fmt.Print` forbidden by pattern `^(fmt.Print(|f|ln)|print|println)$` (forbidigo)\n")
	fmt.Fprintf(os.Stdout, "pkg/lint/lint_test.go:24:23: unused-parameter: parameter 't' seems to be unused, consider removing or renaming it as _ (revive)\n")
	os.Exit(0)
}

func TestGoLint(t *testing.T) {
	utils.ExecCmd = mockGoExec
	defer func() { utils.ExecCmd = exec.Command }() // Restore after test

	// Call golint.Exec
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

	if outputs[0].Path != "pkg/lint/golint.go" {
		t.Errorf("Expected first issue path to be 'pkg/lint/golint.go', got '%s'", outputs[0].Path)
	}

}
