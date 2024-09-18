package utils_test

import (
	"errors"
	"testing"

	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetChangedFiles(t *testing.T) {
	tests := []struct {
		name        string
		mocks       []testutils.MockContent
		expected    []string
		expectedErr error
	}{
		{
			name: "success",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"branch"}},
					Out:    "main\n",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"diff", "--name-only", "main", "--relative"}},
					Out:    "README.md\n",
					Err:    nil,
				},
			},
			expected:    []string{"README.md"},
			expectedErr: nil,
		},
		{
			name: "failure (not in git repository)",
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
			name: "failure (error in git diff command)",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"branch"}},
					Out:    "main\n",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"diff", "--name-only", "main", "--relative"}},
					Out:    "",
					Err:    errors.New("error in git diff command"),
				},
			},
			expected:    []string{},
			expectedErr: errors.New("error in git diff command"),
		},
		{
			name: "failure (error in git branch command)",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"branch"}},
					Out:    "",
					Err:    errors.New("error in git branch command"),
				},
			},
			expected:    []string{},
			expectedErr: errors.New("error in git branch command"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExe := testutils.NewMockExecutor(tt.mocks)
			changedFiles, err := utils.GetChangedFiles(mockExe)
			assert.Equal(t, tt.expected, changedFiles)
			assert.Equal(t, tt.expectedErr, err)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := utils.CheckFilesUpdated(tt.changedFiles, tt.patterns)
			assert.Equal(t, tt.expectedMatch, match)
		})
	}
}
