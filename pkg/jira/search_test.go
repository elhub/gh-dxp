package jira

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTokenFile(t *testing.T, dir, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".jira_token"), []byte(content), 0600))
}

func TestResolveCredentials(t *testing.T) {
	tests := []struct {
		name        string
		configEmail string
		envUsername string
		envToken    string
		fileContent string
		fileExists  bool
		wantEmail   string
		wantToken   string
		wantErr     error
	}{
		{
			name:       "no file, no env vars → not configured",
			fileExists: false,
			wantErr:    ErrJiraNotConfigured,
		},
		{
			name:        "empty file → disabled",
			fileExists:  true,
			fileContent: "",
			wantErr:     ErrJiraDisabled,
		},
		{
			name:        "file contains 'disabled' → disabled",
			fileExists:  true,
			fileContent: "disabled",
			wantErr:     ErrJiraDisabled,
		},
		{
			name:        "file contains email:token → resolved",
			fileExists:  true,
			fileContent: "bot@elhub.no:mytoken",
			wantEmail:   "bot@elhub.no",
			wantToken:   "mytoken",
		},
		{
			name:        "file contains token only → uses config email",
			configEmail: "config@elhub.no",
			fileExists:  true,
			fileContent: "tokenonly",
			wantEmail:   "config@elhub.no",
			wantToken:   "tokenonly",
		},
		{
			name:        "JIRA_API_TOKEN env var overrides file",
			envToken:    "envtoken",
			fileExists:  true,
			fileContent: "file@elhub.no:filetoken",
			wantEmail:   "",
			wantToken:   "envtoken",
		},
		{
			name:        "JIRA_USERNAME env var overrides file email",
			envUsername: "env@elhub.no",
			envToken:    "envtoken",
			fileExists:  false,
			wantEmail:   "env@elhub.no",
			wantToken:   "envtoken",
		},
		{
			name:        "file with whitespace trimmed",
			fileExists:  true,
			fileContent: "  bot@elhub.no : mytoken  ",
			wantEmail:   "bot@elhub.no",
			wantToken:   "mytoken",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			t.Setenv("HOME", tmpDir)
			t.Setenv("JIRA_API_TOKEN", tt.envToken)
			t.Setenv("JIRA_USERNAME", tt.envUsername)

			if tt.fileExists {
				writeTokenFile(t, tmpDir, tt.fileContent)
			}

			email, token, err := resolveCredentials(tt.configEmail)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantEmail, email)
			assert.Equal(t, tt.wantToken, token)
		})
	}
}

func TestRankIssues(t *testing.T) {
	issues := []SearchIssue{
		{Key: "TDX-1", Fields: struct {
			Summary     string          `json:"summary"`
			Description json.RawMessage `json:"description"`
		}{Summary: "fix billing calculation bug"}},
		{Key: "TDX-2", Fields: struct {
			Summary     string          `json:"summary"`
			Description json.RawMessage `json:"description"`
		}{Summary: "add user authentication feature"}},
		{Key: "TDX-3", Fields: struct {
			Summary     string          `json:"summary"`
			Description json.RawMessage `json:"description"`
		}{Summary: "unrelated task"}},
	}

	t.Run("ranks by commit message relevance", func(t *testing.T) {
		result := rankIssues(issues, SearchText{CommitMessage: "fix billing"})
		assert.Equal(t, "TDX-1", result[0].Key)
	})

	t.Run("returns no matches when nothing relevant", func(t *testing.T) {
		result := rankIssues(issues, SearchText{CommitMessage: "xyz123"})
		assert.Empty(t, result)
	})

	t.Run("returns all matching issues", func(t *testing.T) {
		result := rankIssues(issues, SearchText{CommitMessage: "fix billing add user"})
		assert.Len(t, result, 2)
	})
}
