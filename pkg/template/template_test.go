package template_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/template"
	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "template-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Set up the test settings
	settings := &config.Settings{
		ProjectTemplateURI: "http://example.com/template/",
	}

	// Set up the expected file paths
	expectedFilePaths := []string{
		filepath.Join(tempDir, ".editorconfig"),
		filepath.Join(tempDir, ".gitattributes"),
		filepath.Join(tempDir, ".gitignore"),
		filepath.Join(tempDir, "README.md"),
		filepath.Join(tempDir, ".github/CODEOWNERS"),
	}

	// Execute the function
	err = template.Execute(tempDir, settings)
	require.NoError(t, err)

	// Check if the files were created
	for _, filePath := range expectedFilePaths {
		_, err := os.Stat(filePath)
		require.NoError(t, err)
	}

	// Clean up the temporary directory
	os.RemoveAll(tempDir)
}
