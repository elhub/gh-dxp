package lint_test

import (
	"os"
	"testing"

	"github.com/elhub/gh-dxp/pkg/lint"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetekt(t *testing.T) {
	mockExe := new(mocklintExecutor)

	tests := []struct {
		name          string
		mockReturn    string
		mockError     error
		expectedLines int
		expectErr     bool
	}{
		{
			name: "valid detekt output from command",
			mockReturn: "MatchingDeclarationName - [GetPlatformEntriesResponse] at " +
				"/home/ola.norman/work/git/test-app/src/main/kotlin/no/elhub/test-app/data/response/" +
				"GetPlatformDeployLogsResponse.kt:7:1\n" +
				"TopLevelPropertyNaming - [baseRoute] at " +
				"/home/ola.norman/work/git/test-app/src/main/kotlin/no/elhub/test-app/routes/constants.kt:3:11\n" +
				"MagicNumber - [BLUE] at /home/ola.norman/work/git/devxp-deploy-logger/build.gradle.kts:171:10\n",
			mockError:     nil,
			expectedLines: 3,
			expectErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up expectation
			fileString, fileError := lint.GetFiles(",", ".kt")
			require.NoError(t, fileError)
			mockExe.On("Command", "detekt", []string{"-i", fileString, "-r", "md:detekt.out"}).Return(tt.mockReturn, tt.mockError)
			mockExe.On("Command", "rm", []string{"detekt.out"}).Return("", nil)

			// Create a dummy detekt.out file
			writeErr := os.WriteFile("detekt.out", []byte(tt.mockReturn), 0644)
			require.NoError(t, writeErr)

			// Call the method under test
			outputs, err := lint.Detekt{}.Run(mockExe)

			// Clean up detekt.out file
			removeErr := os.Remove("detekt.out")
			require.NoError(t, removeErr)

			// Assert that the expectations were met
			require.Len(t, outputs, tt.expectedLines)
			assert.Equal(t, "detekt", outputs[0].Linter)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDetektParser(t *testing.T) {
	tests := []struct {
		name       string
		inputLine  string
		inputError error
		expected   lint.LinterOutput
		expectErr  bool
	}{
		{
			name: "valid detekt line with line number and column",
			inputLine: "MagicNumber - [<anonymous>] at /home/ola.norman/work/git/test-app/" +
				"src/main/kotlin/no/elhub/test/App.kt:423:4",
			expectErr: false,
			expected: lint.LinterOutput{
				Linter:      "detekt",
				Path:        "/home/ola.norman/work/git/test-app/src/main/kotlin/no/elhub/test/App.kt",
				Line:        423,
				Column:      4,
				Description: "MagicNumber in <anonymous>",
				Severity:    "error",
				Source:      "detekt",
			},
		},
		{
			name: "valid detekt line with line number and column",
			inputLine: "LongMethod - 116/60 - [platformsResourceDocumentation] at /home/ola.norman/work/" +
				"git/test-app/src/main/kotlin/no/elhub/test/App.kt:78:11",
			expectErr: false,
			expected: lint.LinterOutput{
				Linter:      "detekt",
				Path:        "/home/ola.norman/work/git/test-app/src/main/kotlin/no/elhub/test/App.kt",
				Line:        78,
				Column:      11,
				Description: "LongMethod - 116/60 in platformsResourceDocumentation",
				Severity:    "error",
				Source:      "detekt",
			},
		},
		{
			name: "invalid format in detekt message",
			inputLine: "LongMethod - 116/60 - [platformsResourceDocumentation] av /home/ola.norman/work/" +
				"git/test-app/src/main/kotlin/no/elhub/test/App.kt:78:11",
			expectErr: true,
			expected:  lint.LinterOutput{},
		},
		{
			name: "invalid line format in detekt",
			inputLine: "MagicNumber - [<anonymous>] at /home/ola.norman/work/git/test-app/src/main/kotlin/" +
				"no/elhub/test/App.kt:12A:4",
			expectErr: true,
			expected:  lint.LinterOutput{},
		},
		{
			name: "invalid column format in detekt",
			inputLine: "MagicNumber - [<anonymous>] at /home/ola.norman/work/git/test-app/src/main/kotlin/no/" +
				"elhub/test/App.kt:423:X",
			expectErr: true,
			expected:  lint.LinterOutput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the method under test
			output, err := lint.DetektParser(tt.inputLine)

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
