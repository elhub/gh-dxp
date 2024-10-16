package lint_test

import (
	"errors"
	"testing"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/lint"
	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRun_LintNoErrors(t *testing.T) {
	mockExe := new(testutils.MockExecutor)
	mockExe.On("Command", "git", []string{"branch"}).Return("main\ndifferentBranch\n", nil)
	mockExe.On("Command", "git", []string{"diff", "--name-only", "main", "--relative"}).Return("/pkg/source.go\n/pkg/source2.go", nil)

	linterArgs := []string{
		"mega-linter-runner", "--flavor", "cupcake", "-e",
		"MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml",
		"--filesonly", "/pkg/source.go", "/pkg/source2.go",
	}

	mockExe.On("CommandContext", mock.Anything, "npx", linterArgs).Return(nil, nil)

	err := lint.Run(mockExe, &config.Settings{}, &lint.Options{})
	require.NoError(t, err)
	mockExe.AssertExpectations(t)
}

func TestRun_LintHasErrors(t *testing.T) {
	mockExe := new(testutils.MockExecutor)
	mockExe.On("Command", "git", []string{"branch"}).Return("main\ndifferentBranch\n", nil)
	mockExe.On("Command", "git", []string{"diff", "--name-only", "main", "--relative"}).Return("/pkg/source.go\n/pkg/source2.go", nil)

	linterArgs := []string{
		"mega-linter-runner", "--flavor", "cupcake", "-e",
		"MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml",
		"--filesonly", "/pkg/source.go", "/pkg/source2.go",
	}

	mockExe.On("CommandContext", mock.Anything, "npx", linterArgs).Return(nil, errors.New("command error"))

	err := lint.Run(mockExe, &config.Settings{}, &lint.Options{})
	require.Error(t, err)
	mockExe.AssertExpectations(t)
}

func TestRun_LintAllFiles(t *testing.T) {
	mockExe := new(testutils.MockExecutor)

	linterArgs := []string{
		"mega-linter-runner", "--flavor", "cupcake", "-e",
		"MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml",
	}

	mockExe.On("CommandContext", mock.Anything, "npx", linterArgs).Return(nil, nil)

	err := lint.Run(mockExe, &config.Settings{}, &lint.Options{LintAll: true})
	require.NoError(t, err)
	mockExe.AssertExpectations(t)
}

func TestRun_LintWithFix(t *testing.T) {
	mockExe := new(testutils.MockExecutor)
	mockExe.On("Command", "git", []string{"branch"}).Return("main\ndifferentBranch\n", nil)
	mockExe.On("Command", "git", []string{"diff", "--name-only", "main", "--relative"}).Return("/pkg/source.go\n/pkg/source2.go", nil)

	linterArgs := []string{
		"mega-linter-runner", "--flavor", "cupcake", "-e",
		"MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml",
		"--filesonly", "/pkg/source.go", "/pkg/source2.go", "--fix",
	}

	mockExe.On("CommandContext", mock.Anything, "npx", linterArgs).Return(nil, nil)

	err := lint.Run(mockExe, &config.Settings{}, &lint.Options{Fix: true})
	require.NoError(t, err)
	mockExe.AssertExpectations(t)
}

func TestRun_LintWithNoExistingBranches(t *testing.T) {
	mockExe := new(testutils.MockExecutor)
	mockExe.On("Command", "git", []string{"branch"}).Return("", nil)
	mockExe.On("Command", "git", []string{"status", "--porcelain"}).Return(" M /pkg/source.go\n M /pkg/source2.go", nil)

	linterArgs := []string{
		"mega-linter-runner", "--flavor", "cupcake", "-e",
		"MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml",
		"--filesonly", "/pkg/source.go", "/pkg/source2.go",
	}

	mockExe.On("CommandContext", mock.Anything, "npx", linterArgs).Return(nil, nil)

	err := lint.Run(mockExe, &config.Settings{}, &lint.Options{})
	require.NoError(t, err)
	mockExe.AssertExpectations(t)
}

func TestRun_LintSpecificDirectory(t *testing.T) {
	mockExe := new(testutils.MockExecutor)

	linterArgs := []string{
		"mega-linter-runner", "--flavor", "cupcake", "-e",
		"MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml",
		"-e", "FILTER_REGEX_INCLUDE=(pkg)",
	}

	mockExe.On("CommandContext", mock.Anything, "npx", linterArgs).Return(nil, nil)

	err := lint.Run(mockExe, &config.Settings{}, &lint.Options{Directory: "pkg"})
	require.NoError(t, err)
	mockExe.AssertExpectations(t)
}

func TestRun_UseProxy(t *testing.T) {
	mockExe := new(testutils.MockExecutor)

	linterArgs := []string{
		"mega-linter-runner", "--flavor", "cupcake", "-e",
		"MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml",
		"-e", "FILTER_REGEX_INCLUDE=(pkg)", "-e", "https_proxy=https://myproxy.no:8080",
	}

	mockExe.On("CommandContext", mock.Anything, "npx", linterArgs).Return(nil, nil)

	err := lint.Run(mockExe, &config.Settings{}, &lint.Options{Directory: "pkg", Proxy: "https://myproxy.no:8080"})
	require.NoError(t, err)
	mockExe.AssertExpectations(t)
}
