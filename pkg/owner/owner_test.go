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

func TestGetDefaultFile(t *testing.T) {
	tests := []struct {
		name        string
		setupFiles  func(t *testing.T, rootDir string)
		mockContent []testutils.MockContent
		expected    string
		expectErr   bool
		errContains string
	}{
		{
			name: "CODEOWNERS exists in .github",
			setupFiles: func(t *testing.T, rootDir string) {
				err := os.MkdirAll(filepath.Join(rootDir, ".github"), 0755)
				require.NoError(t, err)
				f, err := os.Create(filepath.Join(rootDir, ".github", "CODEOWNERS"))
				require.NoError(t, err)
				f.Close()
			},
			mockContent: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"rev-parse", "--show-toplevel"}},
					Out:    "ROOT_DIR_PLACEHOLDER\n", // Will be replaced in test loop
					Err:    nil,
				},
			},
			expected:  filepath.Join(".github", "CODEOWNERS"), // Relative part we expect to be joined
			expectErr: false,
		},
		{
			name: "CODEOWNERS does not exist",
			setupFiles: func(_ *testing.T, _ string) {
				// Do nothing, file doesn't exist
			},
			mockContent: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"rev-parse", "--show-toplevel"}},
					Out:    "ROOT_DIR_PLACEHOLDER\n",
					Err:    nil,
				},
			},
			expectErr:   true,
			errContains: "could not find CODEOWNERS file",
		},
		{
			name: "Git root error",
			setupFiles: func(_ *testing.T, _ string) {
				// NOOP
			},
			mockContent: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"rev-parse", "--show-toplevel"}},
					Out:    "",
					Err:    errors.New("git error"),
				},
			},
			expectErr:   true,
			errContains: "git error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temp directory to act as the repo root
			tempDir := t.TempDir()

			// Setup files in that temp directory
			if tt.setupFiles != nil {
				tt.setupFiles(t, tempDir)
			}

			// Prepare mocks, replacing the placeholder with the actual temp dir path
			preparedMocks := make([]testutils.MockContent, len(tt.mockContent))
			for i, m := range tt.mockContent {
				if m.Out == "ROOT_DIR_PLACEHOLDER\n" {
					m.Out = tempDir + "\n"
				}
				preparedMocks[i] = m
			}

			mockExe := testutils.NewMockExecutor(preparedMocks)

			// Run the function
			got, err := owner.GetDefaultFile(mockExe)

			if tt.expectErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				// We expect the full path, so join the tempDir with the expected relative path
				expectedPath := filepath.Join(tempDir, tt.expected)
				assert.Equal(t, expectedPath, got)
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
