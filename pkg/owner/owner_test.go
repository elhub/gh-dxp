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
