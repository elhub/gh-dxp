package lint_test

import (
	"context"
	"testing"

	"github.com/elhub/gh-devxp/pkg/config"
	"github.com/elhub/gh-devxp/pkg/lint"
)

// Test main.go
type TestMockLint struct{}

func (TestMockLint) Exec() ([]lint.LintOutput, error) {
	return []lint.LintOutput{
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
	lint.Linters = map[string]lint.Linter{
		"mocklint": TestMockLint{},
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
