package lint_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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
		name             string
		executionError   error
		expectErr        bool
		modifiedFiles    string
		lintAllFiles     bool
		fixFiles         bool
		existingBranches string
		currentChanges   string
		directory        string
	}{
		{
			name:             "lint has no errors",
			executionError:   nil,
			expectErr:        false,
			modifiedFiles:    "/pkg/source.go\n/pkg/source2.go",
			existingBranches: "main\ndifferentBranch\n",
		},
		{
			name:             "lint has errors",
			executionError:   errors.New("command error"),
			expectErr:        true,
			existingBranches: "main\ndifferentBranch\n",
			modifiedFiles:    "/pkg/source.go\n/pkg/source2.go",
		},
		{
			name:             "lint with --all flag",
			executionError:   nil,
			expectErr:        false,
			existingBranches: "main\ndifferentBranch\n",
			lintAllFiles:     true,
		},
		{
			name:             "lint with --fix flag",
			executionError:   nil,
			expectErr:        false,
			modifiedFiles:    "/pkg/source.go\n/pkg/source2.go",
			existingBranches: "main\ndifferentBranch\n",
			fixFiles:         true,
		},
		{
			name:             "lint with no existing branches",
			executionError:   nil,
			expectErr:        false,
			existingBranches: "",
			currentChanges:   " M /pkg/source.go\n M /pkg/source2.go",
			modifiedFiles:    "/pkg/source.go\n/pkg/source2.go",
		},
		{
			name:             "lint the pkg directory",
			executionError:   nil,
			expectErr:        false,
			currentChanges:   " M /pkg/source.go\n M /pkg/source2.go",
			existingBranches: "main\ndifferentBranch\n",
			directory:        "pkg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExe := new(mockExecutor)
			linterArgs := []string{"mega-linter-runner", "--flavor", "cupcake", "-e",
				"MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml"}

			if !tt.lintAllFiles && tt.directory == "" {
				linterArgs = append(linterArgs, "--filesonly")
				linterArgs = append(linterArgs, lint.ConvertTerminalOutputIntoList(tt.modifiedFiles)...)
			} else if tt.directory != "" {
				linterArgs = append(linterArgs, "-e", fmt.Sprintf("FILTER_REGEX_INCLUDE=(%s)", tt.directory))
			}

			if tt.fixFiles {
				linterArgs = append(linterArgs, "--fix")
			}

			mockExe.On("CommandContext", mock.Anything, "npx", linterArgs).Return(nil, tt.executionError)

			if !tt.lintAllFiles && tt.directory == "" {
				mockExe.On("Command", "git", []string{"branch"}).Return(tt.existingBranches, nil)

				if len(tt.existingBranches) == 0 {
					mockExe.On("Command", "git", []string{"status", "--porcelain"}).Return(tt.currentChanges, nil)
				} else {
					mockExe.On("Command", "git", []string{"diff", "--name-only", "main", "--relative"}).Return(tt.modifiedFiles, nil)
				}

			}

			err := lint.Run(mockExe, &config.Settings{}, &lint.Options{LintAll: tt.lintAllFiles, Fix: tt.fixFiles, Directory: tt.directory})

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			mockExe.AssertExpectations(t)
		})
	}
}
