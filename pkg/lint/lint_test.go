package lint_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/lint"
	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRun_LintNoErrors(t *testing.T) {
	testConfig := &config.Settings{
		MegalinterImageVersion: "oxsecurity/megalinter-cupcake:v9",
	}
	mockExe := new(testutils.MockExecutor)
	mockExe.On("Command", "git", []string{"branch"}).Return("main\ndifferentBranch\n", nil)
	mockExe.On("Command", "git", []string{"fetch", "origin", "main"}).Return("", nil)
	mockExe.On("Command", "git", []string{"remote", "set-head", "origin", "--auto"}).Return("", nil)
	mockExe.On("Command", "git", []string{"symbolic-ref", "--short", "refs/remotes/origin/HEAD"}).Return("origin/main", nil)
	mockExe.On("Command", "git", []string{"diff", "--name-only", "origin/main", "--relative"}).Return("/pkg/source.go\n/pkg/source2.go", nil)

	linterArgs := []string{
		"mega-linter-runner", "--image", testConfig.MegalinterImageVersion,
		"-e", "LINTER_RULES_PATH=/tmp",
		"-e", "GOTOOLCHAIN=auto",
		"-e", "MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml",
		"--filesonly", "/pkg/source.go", "/pkg/source2.go",
	}

	mockExe.On("CommandContext", mock.Anything, "npx", linterArgs).Return(nil, nil)

	err := lint.Run(mockExe, testConfig, &lint.Options{})
	require.NoError(t, err)
	mockExe.AssertExpectations(t)
}

func TestRun_LintHasErrors(t *testing.T) {
	testConfig := &config.Settings{
		MegalinterImageVersion: "oxsecurity/megalinter-cupcake:v9",
	}
	mockExe := new(testutils.MockExecutor)
	mockExe.On("Command", "git", []string{"branch"}).Return("main\ndifferentBranch\n", nil)
	mockExe.On("Command", "git", []string{"fetch", "origin", "main"}).Return("", nil)
	mockExe.On("Command", "git", []string{"remote", "set-head", "origin", "--auto"}).Return("", nil)
	mockExe.On("Command", "git", []string{"symbolic-ref", "--short", "refs/remotes/origin/HEAD"}).Return("origin/main", nil)
	mockExe.On("Command", "git", []string{"diff", "--name-only", "origin/main", "--relative"}).Return("/pkg/source.go\n/pkg/source2.go", nil)

	linterArgs := []string{
		"mega-linter-runner", "--image", testConfig.MegalinterImageVersion,
		"-e", "LINTER_RULES_PATH=/tmp",
		"-e", "GOTOOLCHAIN=auto",
		"-e", "MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml",
		"--filesonly", "/pkg/source.go", "/pkg/source2.go",
	}

	mockExe.On("CommandContext", mock.Anything, "npx", linterArgs).Return(nil, errors.New("command error"))

	err := lint.Run(mockExe, testConfig, &lint.Options{})
	require.Error(t, err)
	mockExe.AssertExpectations(t)
}

func TestRun_LintAllFiles(t *testing.T) {
	testConfig := &config.Settings{
		MegalinterImageVersion: "oxsecurity/megalinter-cupcake:v9",
	}
	mockExe := new(testutils.MockExecutor)

	linterArgs := []string{
		"mega-linter-runner", "--image", testConfig.MegalinterImageVersion,
		"-e", "LINTER_RULES_PATH=/tmp",
		"-e", "GOTOOLCHAIN=auto",
		"-e", "MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml",
	}

	mockExe.On("CommandContext", mock.Anything, "npx", linterArgs).Return(nil, nil)

	err := lint.Run(mockExe, testConfig, &lint.Options{LintAll: true})
	require.NoError(t, err)
	mockExe.AssertExpectations(t)
}

func TestRun_LintWithFix(t *testing.T) {
	testConfig := &config.Settings{
		MegalinterImageVersion: "oxsecurity/megalinter-cupcake:v9",
	}
	mockExe := new(testutils.MockExecutor)
	mockExe.On("Command", "git", []string{"branch"}).Return("main\ndifferentBranch\n", nil)
	mockExe.On("Command", "git", []string{"fetch", "origin", "main"}).Return("", nil)
	mockExe.On("Command", "git", []string{"remote", "set-head", "origin", "--auto"}).Return("", nil)
	mockExe.On("Command", "git", []string{"symbolic-ref", "--short", "refs/remotes/origin/HEAD"}).Return("origin/main", nil)
	mockExe.On("Command", "git", []string{"diff", "--name-only", "origin/main", "--relative"}).Return("/pkg/source.go\n/pkg/source2.go", nil)

	linterArgs := []string{
		"mega-linter-runner", "--image", testConfig.MegalinterImageVersion,
		"-e", "LINTER_RULES_PATH=/tmp",
		"-e", "GOTOOLCHAIN=auto",
		"-e", "MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml",
		"--filesonly", "/pkg/source.go", "/pkg/source2.go", "--fix",
	}

	mockExe.On("CommandContext", mock.Anything, "npx", linterArgs).Return(nil, nil)

	err := lint.Run(mockExe, testConfig, &lint.Options{Fix: true})
	require.NoError(t, err)
	mockExe.AssertExpectations(t)
}

func TestRun_LintWithNoExistingBranches(t *testing.T) {
	testConfig := &config.Settings{
		MegalinterImageVersion: "oxsecurity/megalinter-cupcake:v9",
	}
	mockExe := new(testutils.MockExecutor)
	mockExe.On("Command", "git", []string{"branch"}).Return("", nil)
	mockExe.On("Command", "git", []string{"status", "--porcelain"}).Return(" M /pkg/source.go\n M /pkg/source2.go", nil)

	linterArgs := []string{
		"mega-linter-runner", "--image", testConfig.MegalinterImageVersion,
		"-e", "LINTER_RULES_PATH=/tmp",
		"-e", "GOTOOLCHAIN=auto",
		"-e", "MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml",
		"--filesonly", "/pkg/source.go", "/pkg/source2.go",
	}

	mockExe.On("CommandContext", mock.Anything, "npx", linterArgs).Return(nil, nil)

	err := lint.Run(mockExe, testConfig, &lint.Options{})
	require.NoError(t, err)
	mockExe.AssertExpectations(t)
}

func TestRun_LintSpecificDirectory(t *testing.T) {
	testConfig := &config.Settings{
		MegalinterImageVersion: "oxsecurity/megalinter-cupcake:v9",
	}
	mockExe := new(testutils.MockExecutor)

	linterArgs := []string{
		"mega-linter-runner", "--image", testConfig.MegalinterImageVersion,
		"-e", "LINTER_RULES_PATH=/tmp",
		"-e", "GOTOOLCHAIN=auto",
		"-e", "MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml",
		"-e", "FILTER_REGEX_INCLUDE=(pkg)",
	}

	mockExe.On("CommandContext", mock.Anything, "npx", linterArgs).Return(nil, nil)

	err := lint.Run(mockExe, testConfig, &lint.Options{Directory: "pkg"})
	require.NoError(t, err)
	mockExe.AssertExpectations(t)
}

func TestRun_UseProxy(t *testing.T) {
	testConfig := &config.Settings{
		MegalinterImageVersion: "oxsecurity/megalinter-cupcake:v9",
	}
	mockExe := new(testutils.MockExecutor)

	linterArgs := []string{
		"mega-linter-runner", "--image", testConfig.MegalinterImageVersion,
		"-e", "LINTER_RULES_PATH=/tmp",
		"-e", "GOTOOLCHAIN=auto",
		"-e", "MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml",
		"-e", "FILTER_REGEX_INCLUDE=(pkg)", "-e", "https_proxy=https://myproxy.no:8080",
	}

	mockExe.On("CommandContext", mock.Anything, "npx", linterArgs).Return(nil, nil)

	err := lint.Run(mockExe, testConfig, &lint.Options{Directory: "pkg", Proxy: "https://myproxy.no:8080"})
	require.NoError(t, err)
	mockExe.AssertExpectations(t)
}

func TestRun_RepoPinnedGoToolchainSkipsAutoFlag(t *testing.T) {
	testConfig := &config.Settings{
		MegalinterImageVersion: "oxsecurity/megalinter-cupcake:v9",
	}
	mockExe := new(testutils.MockExecutor)

	dir, err := os.MkdirTemp(".", "lint-basedir-")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(dir))
	})

	configContent := "GOTOOLCHAIN: 'go1.27.1'\nGOCACHE: '/tmp/lint/.gocache/build'\n"
	configPath := filepath.Join(dir, ".mega-linter.yml")
	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0o600))

	expectedContainerPath := filepath.ToSlash(filepath.Join("/tmp/lint", filepath.Base(dir), ".mega-linter.yml"))

	linterArgs := []string{
		"mega-linter-runner", "--image", testConfig.MegalinterImageVersion,
		"-e", "LINTER_RULES_PATH=/tmp",
		"-e", "MEGALINTER_CONFIG=" + expectedContainerPath,
	}

	mockExe.On("CommandContext", mock.Anything, "npx", linterArgs).Return(nil, nil)

	err = lint.Run(mockExe, testConfig, &lint.Options{LintAll: true, BaseDir: dir})
	require.NoError(t, err)
	mockExe.AssertExpectations(t)
}

func TestRun_BaseDirWithoutConfigUsesDefaultMegalinterConfig(t *testing.T) {
	testConfig := &config.Settings{
		MegalinterImageVersion: "oxsecurity/megalinter-cupcake:v9",
	}
	mockExe := new(testutils.MockExecutor)

	dir := t.TempDir()

	linterArgs := []string{
		"mega-linter-runner", "--image", testConfig.MegalinterImageVersion,
		"-e", "LINTER_RULES_PATH=/tmp",
		"-e", "GOTOOLCHAIN=auto",
		"-e", "MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml",
	}

	mockExe.On("CommandContext", mock.Anything, "npx", linterArgs).Return(nil, nil)

	err := lint.Run(mockExe, testConfig, &lint.Options{LintAll: true, BaseDir: dir})
	require.NoError(t, err)
	mockExe.AssertExpectations(t)
}

func TestRun_BaseDirOutsideWorkspaceWithLocalConfigReturnsError(t *testing.T) {
	testConfig := &config.Settings{
		MegalinterImageVersion: "oxsecurity/megalinter-cupcake:v9",
	}
	mockExe := new(testutils.MockExecutor)

	dir := t.TempDir()
	configContent := "GOTOOLCHAIN: 'go1.27.1'\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".mega-linter.yml"), []byte(configContent), 0o600))

	err := lint.Run(mockExe, testConfig, &lint.Options{LintAll: true, BaseDir: dir})
	require.Error(t, err)
	require.Contains(t, err.Error(), "outside the mounted workspace")
	mockExe.AssertNotCalled(t, "CommandContext", mock.Anything, mock.Anything, mock.Anything)
}

func TestRepoSetsGoToolchain_Pinned(t *testing.T) {
	dir := t.TempDir()
	configContent := "GOTOOLCHAIN: 'go1.27.1'\nGOCACHE: '/tmp/lint/.gocache/build'\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".mega-linter.yml"), []byte(configContent), 0o600))

	require.True(t, lint.RepoSetsGoToolchain(dir))
}

func TestRepoSetsGoToolchain_NoConfigFile(t *testing.T) {
	dir := t.TempDir()

	// No .mega-linter.yml is written in this temp dir.
	require.False(t, lint.RepoSetsGoToolchain(dir))
}

func TestRepoSetsGoToolchain_ConfigWithoutGoToolchain(t *testing.T) {
	dir := t.TempDir()
	configContent := "ENABLE_LINTERS:\n  - GO_GOLANGCI_LINT\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".mega-linter.yml"), []byte(configContent), 0o600))

	require.False(t, lint.RepoSetsGoToolchain(dir))
}

func TestRepoSetsGoToolchain_MalformedConfig(t *testing.T) {
	dir := t.TempDir()
	configContent := "GOTOOLCHAIN: [unterminated\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".mega-linter.yml"), []byte(configContent), 0o600))

	require.False(t, lint.RepoSetsGoToolchain(dir))
}

func TestRepoSetsGoToolchain_ReadFileError(t *testing.T) {
	dir := t.TempDir()

	// Create a directory with the config file name so os.ReadFile fails with a read error.
	require.NoError(t, os.Mkdir(filepath.Join(dir, ".mega-linter.yml"), 0o700))

	require.False(t, lint.RepoSetsGoToolchain(dir))
}
