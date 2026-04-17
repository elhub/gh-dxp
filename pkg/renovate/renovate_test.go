package renovate_test

import (
	"fmt"
	"testing"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/renovate"
	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const MainBranch = "origin/main"

func TestRenovateRunOnChange(t *testing.T) {
	mockExe := new(testutils.MockExecutor)

	mockExe.On("Command", "git", []string{"branch"}).Return("main\ndifferentBranch\n", nil)
	mockExe.On("Command", "git", []string{"fetch", "origin", "main"}).Return("", nil)
	mockExe.On("Command", "git", []string{"remote", "set-head", "origin", "--auto"}).Return("", nil)
	mockExe.On("Command", "git", []string{"symbolic-ref", "--short", "refs/remotes/origin/HEAD"}).Return("origin/main", nil)
	mockExe.On("Command", "git", []string{"diff", "--name-only", MainBranch, "--relative"}).Return(".github/renovate.json", nil)

	renovateArgs := []string{
		"--package", "renovate@43.78.0", "renovate-config-validator", "--strict",
	}

	mockExe.On("CommandContext", mock.Anything, "npx", renovateArgs).Return(nil, nil)

	err := renovate.Run(mockExe, &config.Settings{}, &renovate.Options{})
	require.NoError(t, err)

	for _, call := range mockExe.Calls {
		fmt.Printf("Called: %s with %v\n", call.Method, call.Arguments)
	}

	mockExe.AssertExpectations(t)
}

func TestRenovateNoRunOnNoChange(t *testing.T) {
	mockExe := new(testutils.MockExecutor)

	mockExe.On("Command", "git", []string{"branch"}).Return("main\ndifferentBranch\n", nil)
	mockExe.On("Command", "git", []string{"fetch", "origin", "main"}).Return("", nil)
	mockExe.On("Command", "git", []string{"remote", "set-head", "origin", "--auto"}).Return("", nil)
	mockExe.On("Command", "git", []string{"symbolic-ref", "--short", "refs/remotes/origin/HEAD"}).Return(MainBranch, nil)
	mockExe.On("Command", "git", []string{"diff", "--name-only", MainBranch, "--relative"}).Return("file1.txt", nil)

	err := renovate.Run(mockExe, &config.Settings{}, &renovate.Options{})

	require.NoError(t, err)

	for _, call := range mockExe.Calls {
		fmt.Printf("Called: %s with %v\n", call.Method, call.Arguments)
	}

	mockExe.AssertNotCalled(t, "CommandContext", mock.Anything, "npx", mock.Anything)
	mockExe.AssertExpectations(t)
}

func TestRenovateRunOnNoChangeWithForce(t *testing.T) {
	mockExe := new(testutils.MockExecutor)

	mockExe.On("Command", "git", []string{"branch"}).Return("main\ndifferentBranch\n", nil)
	mockExe.On("Command", "git", []string{"fetch", "origin", "main"}).Return("", nil)
	mockExe.On("Command", "git", []string{"remote", "set-head", "origin", "--auto"}).Return("", nil)
	mockExe.On("Command", "git", []string{"symbolic-ref", "--short", "refs/remotes/origin/HEAD"}).Return(MainBranch, nil)
	mockExe.On("Command", "git", []string{"diff", "--name-only", MainBranch, "--relative"}).Return("file1.txt", nil)

	renovateArgs := []string{
		"--package", "renovate@43.78.0", "renovate-config-validator", "--strict",
	}


	mockExe.On("CommandContext", mock.Anything, "npx", renovateArgs).Return(nil, nil)

	err := renovate.Run(mockExe, &config.Settings{}, &renovate.Options{Force: true})

	require.NoError(t, err)

	for _, call := range mockExe.Calls {
		fmt.Printf("Called: %s with %v\n", call.Method, call.Arguments)
	}

	mockExe.AssertExpectations(t)
}
