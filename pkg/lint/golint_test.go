package lint_test

import (
	"testing"

	"github.com/elhub/gh-dxp/pkg/lint"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoLint(t *testing.T) {
	mockExe := new(mocklintExecutor)

	tests := []struct {
		name          string
		mockReturn    string
		mockError     error
		expectedLines int
		expectErr     bool
	}{
		{
			name: "valid golint output from command",
			mockReturn: "pkg/lint/golint.go:17:2: use of `fmt.Print` forbidden by pattern " +
				"`^(fmt.Print(|f|ln)|print|println)$` (forbidigo)\n" +
				"pkg/lint/lint_test.go:24:23: unused-parameter: parameter 't' seems " +
				"to be unused, consider removing or renaming it as _ (revive)\n",
			mockError:     nil,
			expectedLines: 2,
			expectErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up expectation
			fileString, fileError := lint.GetFiles(".go", " ")
			require.NoError(t, fileError)
			mockExe.On("Command", "golangci-lint", []string{"run", fileString}).Return(tt.mockReturn, tt.mockError)

			// Call the method under test
			outputs, err := lint.GoLint{}.Run(mockExe)

			// Assert that the expectations were met
			require.Len(t, outputs, tt.expectedLines)
			assert.Equal(t, "golint", outputs[0].Linter)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGoLintParser(t *testing.T) {
	tests := []struct {
		name       string
		inputLine  string
		inputError error
		expected   lint.LinterOutput
		expectErr  bool
	}{
		{
			name:      "valid golint line with line number and column",
			inputLine: "pkg/config/config_test.go:32:5: shadow declaration of \"err\" shadows declaration at line 10 (govet)",
			expectErr: false,
			expected: lint.LinterOutput{
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
			expected: lint.LinterOutput{
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
			expected: lint.LinterOutput{
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
			expected:  lint.LinterOutput{},
		},
		{
			name:      "invalid line format in golint",
			inputLine: "test.yaml:t:8: [error] no new line character at the end of file (new-line-at-end-of-file)",
			expectErr: true,
			expected:  lint.LinterOutput{},
		},
		{
			name:      "invalid column format in golint",
			inputLine: "test.yaml:1:t: [error] no new line character at the end of file (new-line-at-end-of-file)",
			expectErr: true,
			expected:  lint.LinterOutput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the method under test
			output, err := lint.GoLintParser(tt.inputLine)

			// Assert that the expectations were met
			assert.Equal(t, tt.expected, output)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
