package projecttemplate_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/projecttemplate"
	"github.com/stretchr/testify/require"
)

// newTestServer returns an httptest TLS server that responds with dummy content for any path.
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	s := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("dummy template content"))
	}))
	t.Cleanup(s.Close)
	return s
}

func TestExecute(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "template-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Use a local TLS test server instead of hitting the network
	server := newTestServer(t)

	// Set up the test settings
	settings := &config.Settings{
		ProjectTemplateURI: server.URL + "/template/",
	}

	options := &projecttemplate.Options{}

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
	err = projecttemplate.Execute(tempDir, settings, options, server.Client())
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

	// Use a local TLS test server instead of hitting the network
	server := newTestServer(t)

	// Set up the test settings
	settings := &config.Settings{
		ProjectTemplateURI: server.URL + "/template/",
	}

	options := &projecttemplate.Options{
		IsGradleProject: true,
		TestRun:         true,
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
	err = projecttemplate.Execute(tempDir, settings, options, server.Client())
	require.NoError(t, err)

	// Check if the files were created
	for _, filePath := range expectedFilePaths {
		_, err := os.Stat(filePath)
		require.NoError(t, err)
	}

	// Clean up the temporary directory
	os.RemoveAll(tempDir)
}

func TestExecuteDeletesExistingRootFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "template-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Use a local TLS test server instead of hitting the network
	server := newTestServer(t)

	// Create existing CODEOWNERS and CONTRIBUTING.md files at root
	codeownersPath := filepath.Join(tempDir, "CODEOWNERS")
	err = os.WriteFile(codeownersPath, []byte("existing content"), 0644)
	require.NoError(t, err)

	contributorsPath := filepath.Join(tempDir, "CONTRIBUTING.md")
	err = os.WriteFile(contributorsPath, []byte("existing content"), 0644)
	require.NoError(t, err)

	// Set up the test settings
	settings := &config.Settings{
		ProjectTemplateURI: server.URL + "/template/",
	}

	// Create a ConfirmFunc that always confirms deletion
	confirmCalls := []string{}
	confirmFunc := func(prompt string) (bool, error) {
		confirmCalls = append(confirmCalls, prompt)
		return true, nil // Confirm deletion
	}

	options := &projecttemplate.Options{
		CustomAskToConfirmFunc: confirmFunc,
	}

	// Execute the function
	err = projecttemplate.Execute(tempDir, settings, options, server.Client())
	require.NoError(t, err)

	// Verify that confirmation was prompted
	require.Len(t, confirmCalls, 2)

	// Verify that files were deleted
	_, err = os.Stat(codeownersPath)
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))

	_, err = os.Stat(contributorsPath)
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
}

func TestExecuteKeepsFilesWhenUserDeclines(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "template-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Use a local TLS test server instead of hitting the network
	server := newTestServer(t)

	// Create existing CODEOWNERS and CONTRIBUTING.md files at root
	codeownersPath := filepath.Join(tempDir, "CODEOWNERS")
	originalContent := "existing content"
	err = os.WriteFile(codeownersPath, []byte(originalContent), 0644)
	require.NoError(t, err)

	contributingPath := filepath.Join(tempDir, "CONTRIBUTING.md")
	err = os.WriteFile(contributingPath, []byte(originalContent), 0644)
	require.NoError(t, err)

	// Set up the test settings
	settings := &config.Settings{
		ProjectTemplateURI: server.URL + "/template/",
	}

	// Create a ConfirmFunc that always declines deletion
	confirmCalls := []string{}
	confirmFunc := func(prompt string) (bool, error) {
		confirmCalls = append(confirmCalls, prompt)
		return false, nil // Decline deletion
	}

	options := &projecttemplate.Options{
		CustomAskToConfirmFunc: confirmFunc,
	}

	// Execute the function
	err = projecttemplate.Execute(tempDir, settings, options, server.Client())
	require.NoError(t, err)

	// Verify that confirmation was prompted
	require.Len(t, confirmCalls, 2)

	// Verify that files were NOT deleted
	content, err := os.ReadFile(codeownersPath)
	require.NoError(t, err)
	require.Equal(t, originalContent, string(content))

	content, err = os.ReadFile(contributingPath)
	require.NoError(t, err)
	require.Equal(t, originalContent, string(content))
}

func TestExecutePromptMessagesForDeletion(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "template-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Use a local TLS test server instead of hitting the network
	server := newTestServer(t)

	// Create existing CODEOWNERS and CONTRIBUTING.md files at root
	codeownersPath := filepath.Join(tempDir, "CODEOWNERS")
	err = os.WriteFile(codeownersPath, []byte("existing"), 0644)
	require.NoError(t, err)

	// Set up the test settings
	settings := &config.Settings{
		ProjectTemplateURI: server.URL + "/template/",
	}

	// Create a ConfirmFunc that captures the prompts
	var prompts []string
	confirmFunc := func(prompt string) (bool, error) {
		prompts = append(prompts, prompt)
		return false, nil
	}

	options := &projecttemplate.Options{
		CustomAskToConfirmFunc: confirmFunc,
	}

	// Execute the function
	err = projecttemplate.Execute(tempDir, settings, options, server.Client())
	require.NoError(t, err)

	// Verify the correct prompt messages
	require.NotEmpty(t, prompts)
	require.Contains(t, prompts[0], "CODEOWNERS")
	require.Contains(t, prompts[0], "template")
}

func TestExecuteConfirmFuncErrorHandling(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "template-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Use a local TLS test server instead of hitting the network
	server := newTestServer(t)

	// Create existing CODEOWNERS file at root
	codeownersPath := filepath.Join(tempDir, "CODEOWNERS")
	err = os.WriteFile(codeownersPath, []byte("existing"), 0644)
	require.NoError(t, err)

	// Set up the test settings
	settings := &config.Settings{
		ProjectTemplateURI: server.URL + "/template/",
	}

	// Create a ConfirmFunc that returns an error
	confirmFunc := func(_ string) (bool, error) {
		return false, errors.New("confirmation error")
	}

	options := &projecttemplate.Options{
		CustomAskToConfirmFunc: confirmFunc,
	}

	// Execute the function and expect an error
	err = projecttemplate.Execute(tempDir, settings, options, server.Client())
	require.Error(t, err)
	require.Contains(t, err.Error(), "confirmation error")
}
