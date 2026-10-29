package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/elhub/gh-dxp/pkg/cmd"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestOpsLevelCommand(t *testing.T) {
	root := t.TempDir()
	exe := new(testutils.MockExecutor)
	exe.On("Command", "git", []string{"rev-parse", "--show-toplevel"}).Return(root+"\n", nil)
	exe.On("GH", []string{"repo", "view", "--json", "name,nameWithOwner,description,url,primaryLanguage"}).Return(
		`{"name":"example","nameWithOwner":"elhub/example","url":"https://github.com/elhub/example","description":"GitHub description"}`, nil)
	command := cmd.OpsLevelCmd(exe)
	var output bytes.Buffer
	command.SetOut(&output)
	command.SetArgs([]string{
		"--no-input", "--type", "library", "--lifecycle", "development", "--tier", "tier_3",
		"--owner", "team-devxp", "--description", "Explicit description", "--language", "Kotlin",
		"--framework", "Ktor", "--system", "example-system", "--package-name", "no.elhub.example",
	})
	require.NoError(t, command.Execute())
	require.Contains(t, output.String(), "Wrote ")
	content, err := os.ReadFile(filepath.Join(root, "opslevel.yml"))
	require.NoError(t, err)
	var document struct {
		Component map[string]any `yaml:"component"`
	}
	require.NoError(t, yaml.Unmarshal(content, &document))
	for key, value := range map[string]string{
		"type": "library", "lifecycle": "development", "tier": "tier_3", "owner": "team-devxp",
		"description": "Explicit description", "language": "Kotlin", "framework": "Ktor", "system": "example-system",
	} {
		require.Equal(t, value, document.Component[key], key)
	}
	require.Equal(t, []any{"example", "no.elhub.example"}, document.Component["aliases"])
	exe.AssertExpectations(t)
}

func TestOpsLevelCommandRegistered(t *testing.T) {
	root, err := cmd.GenerateCmd(&config.Settings{}, "test")
	require.NoError(t, err)
	command, args, err := root.Find([]string{"opslevel"})
	require.NoError(t, err)
	require.Empty(t, args)
	require.Equal(t, "opslevel", command.Name())
}

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
