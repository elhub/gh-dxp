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

	options := &template.Options{}

	// Set up the expected file paths
	expectedFilePaths := []string{
		filepath.Join(tempDir, ".editorconfig"),
		filepath.Join(tempDir, ".gitattributes"),
		filepath.Join(tempDir, ".gitignore"),
		filepath.Join(tempDir, "README.md"),
		filepath.Join(tempDir, ".github/CODEOWNERS"),
		filepath.Join(tempDir, ".github/CONTRIBUTING.md"),
		filepath.Join(tempDir, ".teamcity/pom.xml"),
		filepath.Join(tempDir, ".teamcity/settings.kts"),
	}

	// Execute the function
	err = template.Execute(tempDir, settings, options)
	require.NoError(t, err)

	// Check if the files were created
	for _, filePath := range expectedFilePaths {
		_, err := os.Stat(filePath)
		require.NoError(t, err)
	}

	// Clean up the temporary directory
	os.RemoveAll(tempDir)
}

func TestExecuteGradle(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "template-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Set up the test settings
	settings := &config.Settings{
		ProjectTemplateURI: "http://example.com/template/",
	}

	options := &template.Options{
		true,
	}

	// Set up the expected file paths
	expectedFilePaths := []string{
		// Standard files
		filepath.Join(tempDir, ".editorconfig"),
		filepath.Join(tempDir, ".gitattributes"),
		filepath.Join(tempDir, ".gitignore"),
		filepath.Join(tempDir, "README.md"),
		filepath.Join(tempDir, ".github/CODEOWNERS"),
		filepath.Join(tempDir, ".github/CONTRIBUTING.md"),
		filepath.Join(tempDir, ".teamcity/pom.xml"),
		filepath.Join(tempDir, ".teamcity/settings.kts"),
		// Gradle specific files
		filepath.Join(tempDir, "gradle/wrapper/gradle-wrapper.jar"),
		filepath.Join(tempDir, "gradle/wrapper/gradle-wrapper.properties"),
		filepath.Join(tempDir, "build.gradle.kts"),
		filepath.Join(tempDir, "gradle.properties"),
		filepath.Join(tempDir, "gradlew"),
		filepath.Join(tempDir, "settings.gradle.kts"),
	}

	// Execute the function
	err = template.Execute(tempDir, settings, options)
	require.NoError(t, err)

	// Check if the files were created
	for _, filePath := range expectedFilePaths {
		_, err := os.Stat(filePath)
		require.NoError(t, err)
	}

	// Clean up the temporary directory
	os.RemoveAll(tempDir)
}
