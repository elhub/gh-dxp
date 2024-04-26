package lint_test

import (
	"testing"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/lint"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mocklintExecutor struct {
	mock.Mock
}

func (m *mocklintExecutor) Command(command string, args ...string) (string, error) {
	argsCalled := m.Called(command, args)
	return argsCalled.String(0), argsCalled.Error(1)
}

type TestMockLint struct{}

func (TestMockLint) Run(_ utils.Executor) ([]lint.LinterOutput, error) {
	return []lint.LinterOutput{
		{
			Linter:      "mocklint",
			Path:        "mock/path",
			Line:        1,
			Column:      1,
			Description: "Ipsum lorem dolor sit amet",
			Severity:    "warning",
			Source:      "mock",
		},
	}, nil
}

func TestRun(t *testing.T) {
	t.Run("should run a linter", func(t *testing.T) {
		mockExe := new(mocklintExecutor)
		var testLinters = map[string]lint.Linter{
			"mocklint": TestMockLint{},
		}
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
		err := lint.Run(mockExe, settings, testLinters)

		// Check if Run returns an error
		require.NoError(t, err)
	})
}
