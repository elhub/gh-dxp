package utils_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinuxExecutor_Command(t *testing.T) {
	t.Run("should return command output", func(t *testing.T) {
		executor := &utils.LinuxExecutorImpl{
			ExecCommand: func(_ string, _ ...string) *exec.Cmd {
				return exec.Command("echo", "hello")
			},
		}

		// Call the method under test
		output, err := executor.Command("echo", "hello")

		// Assert that the expectations were met
		require.NoError(t, err)
		assert.Equal(t, "hello\n", output)
	})

	t.Run("should return error if command fails", func(t *testing.T) {
		executor := &utils.LinuxExecutorImpl{
			ExecCommand: func(_ string, _ ...string) *exec.Cmd {
				return exec.Command("ls", "/nonexistent")
			},
		}

		// Call the method under test
		output, err := executor.Command("ls", "/nonexistent")

		// Assert that the expectations were met
		require.Error(t, err)
		assert.Contains(t, output, "No such file or directory")
	})
}

func TestLinuxExecutor_CommandContext(t *testing.T) {
	// Note: These tests rely on 'true' and 'sleep' being available in the path,
	// which is standard for Linux environments.

	t.Run("should execute command successfully", func(t *testing.T) {
		executor := utils.LinuxExecutor()
		ctx := context.Background()

		// 'true' command exits with 0 immediately
		err := executor.CommandContext(ctx, "true")

		require.NoError(t, err)
	})

	t.Run("should return error for failing command", func(t *testing.T) {
		executor := utils.LinuxExecutor()
		ctx := context.Background()

		// 'false' command exits with 1
		err := executor.CommandContext(ctx, "false")

		require.Error(t, err)
	})

	t.Run("should kill process on context cancellation", func(t *testing.T) {
		executor := utils.LinuxExecutor()

		// Create a context that cancels quickly (100ms)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		// Run a command that sleeps for longer than the timeout (2s)
		start := time.Now()
		err := executor.CommandContext(ctx, "sleep", "2")
		duration := time.Since(start)

		// Verify that it returned with a context error
		require.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)

		// Verify that it returned quickly and didn't wait for the full sleep
		assert.Less(t, duration, 1500*time.Millisecond, "Command should have been cancelled immediately")
	})
}

func TestLinuxExecutor_Chdir(t *testing.T) {
	executor := utils.LinuxExecutor()

	// Create a temporary directory for the test
	tmpDir := t.TempDir()

	// Save current working directory to restore it after test
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(originalWd) }()

	// Perform the change
	err = executor.Chdir(tmpDir)
	require.NoError(t, err)

	// Verify the change happened
	currentWd, err := os.Getwd()
	require.NoError(t, err)

	// Evaluate symlinks to ensure paths match (common issue with /tmp on some OSs)
	evalTmpDir, _ := filepath.EvalSymlinks(tmpDir)
	evalCurrentWd, _ := filepath.EvalSymlinks(currentWd)

	assert.Equal(t, evalTmpDir, evalCurrentWd)
}
