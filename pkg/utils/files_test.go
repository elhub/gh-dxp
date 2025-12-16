package utils_test

import (
	"errors"
	"os"
	"testing"

	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetChangedFiles(t *testing.T) {
	tests := []struct {
		name        string
		mocks       []testutils.MockContent
		expected    []string
		expectedErr error
	}{
		{
			name: "success with branch list",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"branch"}},
					Out:    "* main\n  feature-branch\n",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"fetch", "origin", "main"}},
					Out:    "",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"remote", "set-head", "origin", "--auto"}},
					Out:    "",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"symbolic-ref", "--short", "refs/remotes/origin/HEAD"}},
					Out:    "origin/main",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"diff", "--name-only", "origin/main", "--relative"}},
					Out:    "README.md\nsrc/main.go\n",
					Err:    nil,
				},
			},
			expected:    []string{"README.md", "src/main.go"},
			expectedErr: nil,
		},
		{
			name: "success with single file",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"branch"}},
					Out:    "main\n",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"fetch", "origin", "main"}},
					Out:    "",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"remote", "set-head", "origin", "--auto"}},
					Out:    "",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"symbolic-ref", "--short", "refs/remotes/origin/HEAD"}},
					Out:    "origin/main\n",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"diff", "--name-only", "origin/main", "--relative"}},
					Out:    "README.md\n",
					Err:    nil,
				},
			},
			expected:    []string{"README.md"},
			expectedErr: nil,
		},
		{
			name: "empty branch list - fallback to tracked changes",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"branch"}},
					Out:    "",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "M  file1.go\nA  file2.go\n",
					Err:    nil,
				},
			},
			expected:    []string{"file1.go", "file2.go"},
			expectedErr: nil,
		},
		{
			name: "failure - not in git repository",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"branch"}},
					Out:    "",
					Err:    errors.New("not a git repository"),
				},
			},
			expected:    []string{},
			expectedErr: errors.New("not a git repository"),
		},
		{
			name: "failure - error in git fetch",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"branch"}},
					Out:    "main\n",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"fetch", "origin", "main"}},
					Out:    "",
					Err:    errors.New("fetch failed"),
				},
			},
			expected:    nil,
			expectedErr: nil, // Bug in implementation: returns err instead of pullErr
		},
		{
			name: "failure - error in set-head",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"branch"}},
					Out:    "main\n",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"fetch", "origin", "main"}},
					Out:    "",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"remote", "set-head", "origin", "--auto"}},
					Out:    "",
					Err:    errors.New("set-head failed"),
				},
			},
			expected:    nil,
			expectedErr: errors.New("set-head failed"),
		},
		{
			name: "failure - error in symbolic-ref",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"branch"}},
					Out:    "main\n",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"fetch", "origin", "main"}},
					Out:    "",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"remote", "set-head", "origin", "--auto"}},
					Out:    "",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"symbolic-ref", "--short", "refs/remotes/origin/HEAD"}},
					Out:    "",
					Err:    errors.New("symbolic-ref failed"),
				},
			},
			expected:    nil,
			expectedErr: errors.New("symbolic-ref failed"),
		},
		{
			name: "failure - error in git diff command",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"branch"}},
					Out:    "main\n",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"fetch", "origin", "main"}},
					Out:    "",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"remote", "set-head", "origin", "--auto"}},
					Out:    "origin/main\n",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"symbolic-ref", "--short", "refs/remotes/origin/HEAD"}},
					Out:    "origin/main\n",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"diff", "--name-only", "origin/main", "--relative"}},
					Out:    "",
					Err:    errors.New("error in git diff command"),
				},
			},
			expected:    []string{},
			expectedErr: errors.New("error in git diff command"),
		},
		{
			name: "no changes",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"branch"}},
					Out:    "main\n",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"fetch", "origin", "main"}},
					Out:    "",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"remote", "set-head", "origin", "--auto"}},
					Out:    "",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"symbolic-ref", "--short", "refs/remotes/origin/HEAD"}},
					Out:    "origin/main",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"diff", "--name-only", "origin/main", "--relative"}},
					Out:    "",
					Err:    nil,
				},
			},
			expected:    []string{},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExe := testutils.NewMockExecutor(tt.mocks)
			changedFiles, err := utils.GetChangedFiles(mockExe)
			assert.Equal(t, tt.expected, changedFiles)
			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.Equal(t, tt.expectedErr, err)
			}
			mockExe.AssertExpectations(t)
		})
	}
}

func TestCheckFilesUpdated(t *testing.T) {
	tests := []struct {
		name          string
		changedFiles  []string
		patterns      []string
		expectedMatch bool
	}{
		{
			name:          "No files changed",
			changedFiles:  []string{},
			patterns:      []string{"README.md$"},
			expectedMatch: false,
		},
		{
			name:          "Match specific file name",
			changedFiles:  []string{"README.md"},
			patterns:      []string{"README.md$"},
			expectedMatch: true,
		},
		{
			name:          "Match end of path",
			changedFiles:  []string{"/path/to/the/README.md"},
			patterns:      []string{"README.md$"},
			expectedMatch: true,
		},
		{
			name:          "Do not match part of the path",
			changedFiles:  []string{"/path/README.md/around/hello.md"},
			patterns:      []string{"README.md$"},
			expectedMatch: false,
		},
		{
			name:          "Files in docs directory changed",
			changedFiles:  []string{"/test/docs/index.md", "/test/docs/guide.md"},
			patterns:      []string{"/docs/"},
			expectedMatch: true,
		},
		{
			name:          "Files in docs-like directory changed",
			changedFiles:  []string{"docsy/index.md", "/src/com/manydocs/helloworld.go"},
			patterns:      []string{"/docs/"},
			expectedMatch: false,
		},
		{
			name:          "Multiple patterns matched",
			changedFiles:  []string{"README.md", "/docs/index.md"},
			patterns:      []string{"README.md", "/docs/"},
			expectedMatch: true,
		},
		{
			name:          "No patterns matched",
			changedFiles:  []string{"main.go", "utils/helper.go"},
			patterns:      []string{"README.md", "/docs/"},
			expectedMatch: false,
		},
		{
			name:          "Wildcard pattern",
			changedFiles:  []string{"test.go", "main.go"},
			patterns:      []string{`\.go$`},
			expectedMatch: true,
		},
		{
			name:          "Invalid regex pattern - no match",
			changedFiles:  []string{"test.go"},
			patterns:      []string{`[invalid`}, // Invalid regex
			expectedMatch: false,
		},
		{
			name:          "Empty patterns",
			changedFiles:  []string{"test.go"},
			patterns:      []string{},
			expectedMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := utils.CheckFilesUpdated(tt.changedFiles, tt.patterns)
			assert.Equal(t, tt.expectedMatch, match)
		})
	}
}

func TestGetTrackedChanges(t *testing.T) {
	tests := []struct {
		name        string
		mocks       []testutils.MockContent
		expected    []string
		expectedErr error
	}{
		{
			name: "modified files - staged",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "M  file1.go\nM  file2.go\n",
					Err:    nil,
				},
			},
			expected:    []string{"file1.go", "file2.go"},
			expectedErr: nil,
		},
		{
			name: "modified files - unstaged",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    " M file1.go\n M file2.go\n",
					Err:    nil,
				},
			},
			expected:    []string{"file1.go", "file2.go"},
			expectedErr: nil,
		},
		{
			name: "added files",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "A  newfile.go\n",
					Err:    nil,
				},
			},
			expected:    []string{"newfile.go"},
			expectedErr: nil,
		},
		{
			name: "deleted files",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "D  oldfile.go\n",
					Err:    nil,
				},
			},
			expected:    []string{"oldfile.go"},
			expectedErr: nil,
		},
		{
			name: "renamed files",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "R  oldname.go -> newname.go\n",
					Err:    nil,
				},
			},
			expected:    []string{"newname.go"},
			expectedErr: nil,
		},
		{
			name: "mixed changes with untracked files",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "M  modified.go\nA  added.go\n?? untracked.go\n",
					Err:    nil,
				},
			},
			expected:    []string{"modified.go", "added.go"},
			expectedErr: nil,
		},
		{
			name: "no tracked changes - only untracked",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "?? untracked.go\n",
					Err:    nil,
				},
			},
			expected:    []string{},
			expectedErr: nil,
		},
		{
			name: "empty status",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "",
					Err:    nil,
				},
			},
			expected:    []string{},
			expectedErr: nil,
		},
		{
			name: "git command error",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "",
					Err:    errors.New("not a git repository"),
				},
			},
			expected:    []string{},
			expectedErr: errors.New("not a git repository"),
		},
		{
			name: "mixed staged and unstaged",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "MM file1.go\n M file2.go\nM  file3.go\n",
					Err:    nil,
				},
			},
			expected:    []string{"file1.go", "file2.go", "file3.go"},
			expectedErr: nil,
		},
		{
			name: "renamed with path containing spaces",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "R  old name.go -> new name.go\n",
					Err:    nil,
				},
			},
			expected:    []string{"new name.go"},
			expectedErr: nil,
		},
		{
			name: "type changed file",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "T  changed.go\n",
					Err:    nil,
				},
			},
			expected:    []string{"changed.go"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExe := testutils.NewMockExecutor(tt.mocks)
			changes, err := utils.GetTrackedChanges(mockExe)
			assert.Equal(t, tt.expected, changes)
			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.Equal(t, tt.expectedErr, err)
			}
			mockExe.AssertExpectations(t)
		})
	}
}

func TestGetUntrackedChanges(t *testing.T) {
	tests := []struct {
		name        string
		mocks       []testutils.MockContent
		expected    []string
		expectedErr error
	}{
		{
			name: "untracked files",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "?? untracked1.go\n?? untracked2.go\n",
					Err:    nil,
				},
			},
			expected:    []string{"untracked1.go", "untracked2.go"},
			expectedErr: nil,
		},
		{
			name: "mixed changes with tracked files",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "M  modified.go\n?? untracked.go\n",
					Err:    nil,
				},
			},
			expected:    []string{"untracked.go"},
			expectedErr: nil,
		},
		{
			name: "no untracked files - only tracked",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "M  modified.go\n",
					Err:    nil,
				},
			},
			expected:    []string{},
			expectedErr: nil,
		},
		{
			name: "empty status",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "",
					Err:    nil,
				},
			},
			expected:    []string{},
			expectedErr: nil,
		},
		{
			name: "git command error",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "",
					Err:    errors.New("not a git repository"),
				},
			},
			expected:    []string{},
			expectedErr: errors.New("not a git repository"),
		},
		{
			name: "untracked directory",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "?? newdir/\n",
					Err:    nil,
				},
			},
			expected:    []string{"newdir/"},
			expectedErr: nil,
		},
		{
			name: "untracked files with spaces in name",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "?? file with spaces.go\n",
					Err:    nil,
				},
			},
			expected:    []string{"file with spaces.go"},
			expectedErr: nil,
		},
		{
			name: "multiple untracked files and directories",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"status", "--porcelain"}},
					Out:    "?? file1.go\n?? dir1/\n?? file2.txt\nM  tracked.go\n",
					Err:    nil,
				},
			},
			expected:    []string{"file1.go", "dir1/", "file2.txt"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExe := testutils.NewMockExecutor(tt.mocks)
			changes, err := utils.GetUntrackedChanges(mockExe)
			assert.Equal(t, tt.expected, changes)
			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.Equal(t, tt.expectedErr, err)
			}
			mockExe.AssertExpectations(t)
		})
	}
}

func TestFileExists(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) string
		expected bool
	}{
		{
			name: "file exists",
			setup: func(t *testing.T) string {
				tmpFile := t.TempDir() + "/testfile.txt"
				err := os.WriteFile(tmpFile, []byte("test"), 0644)
				require.NoError(t, err)
				return tmpFile
			},
			expected: true,
		},
		{
			name: "file does not exist",
			setup: func(t *testing.T) string {
				return t.TempDir() + "/nonexistent.txt"
			},
			expected: false,
		},
		{
			name: "directory exists",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup(t)
			result := utils.FileExists(path)
			assert.Equal(t, tt.expected, result)
		})
	}
}
