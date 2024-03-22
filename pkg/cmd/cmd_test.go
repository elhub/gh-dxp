package cmd_test

import (
	"os"
	"testing"

	"github.com/elhub/gh-dxp/pkg/cmd"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	t.Run("should execute main command without error", func(t *testing.T) {
		settings := &config.Settings{}
		version := "1.0.0"

		// Redirect output to prevent printing to console during test
		oldOut := os.Stdout
		_, w, _ := os.Pipe()
		os.Stdout = w

		// Call the function under test
		err := cmd.Execute(settings, version)

		// Restore original stdout
		os.Stdout = oldOut

		// Assert that the command executed without error
		// Note: You might need to add more assertions based on the behavior of your command
		require.NoError(t, err)
	})
}
