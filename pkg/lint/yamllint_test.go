package lint_test

import (
	"testing"

	"github.com/elhub/gh-dxp/pkg/lint"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYamlLint(t *testing.T) {
	mockExe := new(mocklintExecutor)

	tests := []struct {
		name          string
		mockReturn    string
		mockError     error
		expectedLines int
		expectErr     bool
	}{
		{
			name: "valid yamllint output from command",
			mockReturn: "./test/weird.yml:1:1: [warning] missing document start \"---\" (document-start)\n" +
				"./test/weird.yml:2:1: [warning] truthy value should be one of [false, true] (truthy)\n",
			mockError:     nil,
			expectedLines: 2,
			expectErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up expectation
			mockExe.On("Command", "git", []string{"rev-parse", "--show-toplevel"}).Return("./resources", nil)
			files, fileError := lint.GetFiles(mockExe, " ", ".yaml", ".yml")
			require.NoError(t, fileError)
			mockExe.On("Command", "yamllint", []string{files, "-f", "parsable", "."}).Return(tt.mockReturn, tt.mockError)

			// Call the method under test
			outputs, err := lint.YamlLint{}.Run(mockExe)

			// Assert that the expectations were met
			require.Len(t, outputs, tt.expectedLines)
			assert.Equal(t, "yamllint", outputs[0].Linter)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestYamlLintParser(t *testing.T) {
	tests := []struct {
		name      string
		inputLine string
		expected  lint.LinterOutput
		expectErr bool
	}{
		{
			name: "valid yammlint line with parenthesis",
			inputLine: "test_infra/test_plugin/new_test/fix.yaml:195:101: [error] line too long (399 > " +
				"100 characters) (line-length)",
			expectErr: false,
			expected: lint.LinterOutput{
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
			expected: lint.LinterOutput{
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
			expected:  lint.LinterOutput{},
		},
		{
			name:      "yamllint invalid column",
			inputLine: "test.yaml:1:t: [error] no new line character at the end of file (new-line-at-end-of-file)",
			expectErr: true,
			expected:  lint.LinterOutput{},
		},
		{
			name:      "yamllint invalid format",
			inputLine: "[error] test:1,2 no new line character at the end of file (new-line-at-end-of-file)",
			expectErr: true,
			expected:  lint.LinterOutput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the method under test
			output, err := lint.YamlLintParser(tt.inputLine)

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
