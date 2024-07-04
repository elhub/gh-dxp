package lint_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/lint"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockExecutor struct {
	mock.Mock
}

func (m *mockExecutor) Command(name string, arg ...string) (string, error) {
	args := m.Called(name, arg)
	return args.String(0), args.Error(1)
}

func (m *mockExecutor) CommandContext(ctx context.Context, name string, arg ...string) error {
	args := m.Called(ctx, name, arg)
	return args.Error(1)
}

func (m *mockExecutor) GH(arg ...string) (bytes.Buffer, error) {
	args := m.Called(arg)
	return *bytes.NewBufferString(args.String(0)), args.Error(1)
}

func TestRun(t *testing.T) {
	tests := []struct {
		name           string
		executionError error
		expectErr      bool
		modifiedFiles  string
		lintAllFiles   bool
	}{
		{
			name:           "lint has no errors",
			executionError: nil,
			expectErr:      false,
			modifiedFiles:  "/pkg/source.go\n/pkg/source2.go",
		},
		{
			name:           "lint has errors",
			executionError: errors.New("command error"),
			expectErr:      true,
			modifiedFiles:  "/pkg/source.go\n/pkg/source2.go",
		},
		{
			name:           "lint with --all flag",
			executionError: nil,
			expectErr:      false,
			lintAllFiles:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExe := new(mockExecutor)
			linterArgs := []string{"mega-linter-runner", "--flavor", "cupcake", "-e",
				"MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml"}

			if !tt.lintAllFiles {
				linterArgs = append(linterArgs, "--filesonly")
				linterArgs = append(linterArgs, lint.ConvertChangedFilesIntoList(tt.modifiedFiles)...)
			}

			mockExe.On("CommandContext", mock.Anything, "npx", linterArgs).Return(nil, tt.executionError)

			if !tt.lintAllFiles {
				mockExe.On("Command", "git", []string{"diff", "--name-only", "main", "--relative"}).Return(tt.modifiedFiles, nil)
			}

			err := lint.Run(mockExe, &config.Settings{}, &lint.Options{LintAll: tt.lintAllFiles})

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			mockExe.AssertExpectations(t)
		})
	}
}
