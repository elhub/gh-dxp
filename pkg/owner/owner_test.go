package owner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/elhub/gh-dxp/pkg/owner"
	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		gitRoot        string
		gitRootErr     error
		codeowners     string
		expectedOwners []string
		expectErr      bool
		errContains    string
	}{
		{
			name:           "Single owner for README",
			path:           "README.md",
			gitRoot:        createTempRepoWithCodeowners(t, `* @elhub/devxp`),
			expectedOwners: []string{"@elhub/devxp"},
			expectErr:      false,
		},
		{
			name: "Multiple owners for specific path",
			path: "src/main.go",
			gitRoot: createTempRepoWithCodeowners(t, `*.go @elhub/backend @elhub/devxp
src/ @elhub/developers`),
			expectedOwners: []string{"@elhub/developers"},
			expectErr:      false,
		},
		{
			name:        "Git root directory error",
			path:        "README.md",
			gitRootErr:  errors.New("not a git repository"),
			expectErr:   true,
			errContains: "not a git repository",
		},
		{
			name:        "CODEOWNERS file not found",
			path:        "README.md",
			gitRoot:     createTempRepoWithoutCodeowners(t),
			expectErr:   true,
			errContains: "no such file or directory",
		},
		{
			name:           "Multiple team owners",
			path:           "pkg/api/handler.go",
			gitRoot:        createTempRepoWithCodeowners(t, `pkg/api/ @elhub/api-team @elhub/backend-team @elhub/devxp`),
			expectedOwners: []string{"@elhub/api-team", "@elhub/backend-team", "@elhub/devxp"},
			expectErr:      false,
		},
		{
			name:           "User and team owners",
			path:           "security/auth.go",
			gitRoot:        createTempRepoWithCodeowners(t, `security/ @user1 @elhub/security-team`),
			expectedOwners: []string{"@user1", "@elhub/security-team"},
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExe := new(testutils.MockExecutor)

			if tt.gitRoot != "" {
				mockExe.On("Command", "git", []string{"rev-parse", "--show-toplevel"}).
					Return(tt.gitRoot, tt.gitRootErr)
			} else {
				mockExe.On("Command", "git", []string{"rev-parse", "--show-toplevel"}).
					Return("", tt.gitRootErr)
			}

			owners, err := owner.Execute(tt.path, mockExe)

			if tt.expectErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedOwners, owners)
			}

			mockExe.AssertExpectations(t)
		})
	}
}

/*
func TestGetDefaultFile(t *testing.T) {
	tests := []struct {
		name         string
		gitRoot      string
		gitRootErr   error
		createFile   bool
		expectedPath string
		expectErr    bool
		errContains  string
	}{
		{
			name:       "CODEOWNERS file exists",
			gitRoot:    createTempRepoWithCodeowners(t, `* @elhub/devxp`),
			createFile: true,
			expectErr:  false,
		},
		{
			name:        "CODEOWNERS file does not exist",
			gitRoot:     createTempRepoWithoutCodeowners(t),
			createFile:  false,
			expectErr:   true,
			errContains: "could not find CODEOWNERS file in .github directory",
		},
		{
			name:        "Git root directory error",
			gitRootErr:  errors.New("not a git repository"),
			expectErr:   true,
			errContains: "not a git repository",
		},
		{
			name:        "Empty git root",
			gitRoot:     "",
			gitRootErr:  errors.New("failed to get root directory"),
			expectErr:   true,
			errContains: "failed to get root directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExe := new(testutils.MockExecutor)

			mockExe.On("Command", "git", []string{"rev-parse", "--show-toplevel"}).
				Return(tt.gitRoot, tt.gitRootErr)

			filePath, err := owner.GetDefaultFile(mockExe)

			if tt.expectErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				expectedPath := filepath.Join(tt.gitRoot, ".github/CODEOWNERS")
				assert.Equal(t, expectedPath, filePath)
				// Verify the file actually exists
				assert.FileExists(t, filePath)
			}

			mockExe.AssertExpectations(t)
		})
	}
}

func TestExecute_InvalidCodeownersFormat(t *testing.T) {
	// Test with malformed CODEOWNERS file
	gitRoot := createTempRepoWithCodeowners(t, `this is not a valid format @@@`)

	mockExe := new(testutils.MockExecutor)
	mockExe.On("Command", "git", []string{"rev-parse", "--show-toplevel"}).
		Return(gitRoot, nil)

	owners, err := owner.Execute("README.md", mockExe)

	// The codeowners library should handle this - check behavior
	// This test documents the actual behavior
	if err != nil {
		require.Error(t, err)
	} else {
		// If no error, owners should still be returned (possibly empty or default)
		assert.NotNil(t, owners)
	}

	mockExe.AssertExpectations(t)
}

func TestExecute_EmptyCodeownersFile(t *testing.T) {
	// Test with empty CODEOWNERS file
	gitRoot := createTempRepoWithCodeowners(t, ``)

	mockExe := new(testutils.MockExecutor)
	mockExe.On("Command", "git", []string{"rev-parse", "--show-toplevel"}).
		Return(gitRoot, nil)

	owners, err := owner.Execute("README.md", mockExe)

	// With an empty CODEOWNERS file, there should be an error or no owners
	if err != nil {
		require.Error(t, err)
	} else {
		assert.Empty(t, owners)
	}

	mockExe.AssertExpectations(t)
}

func TestExecute_ComplexPatterns(t *testing.T) {
	tests := []struct {
		name              string
		path              string
		codeownersContent string
		expectedOwners    []string
	}{
		{
			name: "Pattern with negation and override",
			path: "src/app.js",
			codeownersContent: `# Default owner
* @elhub/default-team

# JavaScript files
*.js @elhub/frontend-team

# Source directory
src/ @elhub/src-team`,
			expectedOwners: []string{"@elhub/src-team"},
		},
		{
			name: "Extension-based matching",
			path: "config.yml",
			codeownersContent: `*.yml @elhub/config-team
*.yaml @elhub/config-team
*.json @elhub/config-team`,
			expectedOwners: []string{"@elhub/config-team"},
		},
		{
			name: "Directory with trailing slash",
			path: "tests/unit/test.go",
			codeownersContent: `tests/ @elhub/qa-team
*.go @elhub/backend-team`,
			expectedOwners: []string{"@elhub/qa-team"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRoot := createTempRepoWithCodeowners(t, tt.codeownersContent)

			mockExe := new(testutils.MockExecutor)
			mockExe.On("Command", "git", []string{"rev-parse", "--show-toplevel"}).
				Return(gitRoot, nil)

			owners, err := owner.Execute(tt.path, mockExe)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedOwners, owners)

			mockExe.AssertExpectations(t)
		})
	}
}
*/
// Helper functions

// createTempRepoWithCodeowners creates a temporary directory structure with .github/CODEOWNERS
func createTempRepoWithCodeowners(t *testing.T, codeownersContent string) string {
	t.Helper()

	tmpDir := t.TempDir()
	githubDir := filepath.Join(tmpDir, ".github")

	err := os.MkdirAll(githubDir, 0755)
	require.NoError(t, err)

	codeownersPath := filepath.Join(githubDir, "CODEOWNERS")
	err = os.WriteFile(codeownersPath, []byte(codeownersContent), 0644)
	require.NoError(t, err)

	return tmpDir
}

// createTempRepoWithoutCodeowners creates a temporary directory without CODEOWNERS file
func createTempRepoWithoutCodeowners(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	// Create .github directory but no CODEOWNERS file
	githubDir := filepath.Join(tmpDir, ".github")
	err := os.MkdirAll(githubDir, 0755)
	require.NoError(t, err)

	return tmpDir
}
