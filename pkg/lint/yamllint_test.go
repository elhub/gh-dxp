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

func TestYamlLintParser(t *testing.T) {
	tests := []struct {
		name      string
		inputLine string
		expected  lint.LintOutput
		expectErr bool
	}{
		{
			name: "valid yammlint line with parenthesis",
			inputLine: "test_infra/test_plugin/new_test/fix.yaml:195:101: [error] line too long (399 > " +
				"100 characters) (line-length)",
			expectErr: false,
			expected: lint.LintOutput{
				Linter:      "yamllint",
				Path:        "test_infra/test_plugin/new_test/fix.yaml",
				Line:        195,
				Column:      101,
				Description: "line too long (399 > 100 characters)",
				Severity:    "error",
				Source:      "line-length",
			},
		},
		{
			name:      "valid yammlint line without parenthesis",
			inputLine: "test.yaml:1:8: [error] no new line character at the end of file (new-line-at-end-of-file)",
			expectErr: false,
			expected: lint.LintOutput{
				Linter:      "yamllint",
				Path:        "test.yaml",
				Line:        1,
				Column:      8,
				Severity:    "error",
				Description: "no new line character at the end of file",
				Source:      "new-line-at-end-of-file",
			},
		},
		{
			name:      "yamllint invalid line",
			inputLine: "test.yaml:t:8: [error] no new line character at the end of file (new-line-at-end-of-file)",
			expectErr: true,
			expected:  lint.LintOutput{},
		},
		{
			name:      "yamllint invalid column",
			inputLine: "test.yaml:1:t: [error] no new line character at the end of file (new-line-at-end-of-file)",
			expectErr: true,
			expected:  lint.LintOutput{},
		},
		{
			name:      "yamllint invalid format",
			inputLine: "[error] test:1,2 no new line character at the end of file (new-line-at-end-of-file)",
			expectErr: true,
			expected:  lint.LintOutput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := lint.YamlLintParser(tt.inputLine)
			if (err != nil) != tt.expectErr {
				t.Errorf("YamlLintParser() output = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !reflect.DeepEqual(output, tt.expected) {
				t.Errorf("YamlLintParser() output = %v, expected %v", output, tt.expected)
			}
		})
	}

}

// Test golint.go
func mockYamlExec(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestYamlLintErrors", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestYamlLintErrors(_ *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintf(os.Stdout, "./test/weird.yml:1:1: [warning] missing document start \"---\" (document-start)\n")
	fmt.Fprintf(os.Stdout, "./test/weird.yml:2:1: [warning] truthy value should be one of [false, true] (truthy)\n")
	os.Exit(0)
}

func TestYamlLint(t *testing.T) {
	utils.ExecCmd = mockYamlExec
	defer func() { utils.ExecCmd = exec.Command }() // Restore after test

	// Call YamlLint.Exec
	outputs, err := lint.YamlLint{}.Exec()

	if err != nil {
		t.Errorf("YamlLint threw an error when none was expected")
	}

	if len(outputs) != 2 {
		t.Errorf("Expected 2 issues in YamlLint.Exec, got '%d'", len(outputs))
	}

	if outputs[0].Linter != "yamllint" {
		t.Errorf("Expected first issue linter to be 'yamllint', got '%s'", outputs[0].Linter)
	}

	if outputs[0].Path != "./test/weird.yml" {
		t.Errorf("Expected first issue path to be './test/weird.yml', got '%s'", outputs[0].Path)
	}

}
