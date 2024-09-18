package utils_test

import (
	"errors"
	"testing"

	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetGitRootDirectory(t *testing.T) {
	tests := []struct {
		name     string
		mocks    []testutils.MockContent
		expected string
		wantErr  bool
	}{
		{
			name: "success",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"rev-parse", "--show-toplevel"}},
					Out:    "/path/to/repo\n",
					Err:    nil,
				},
			},
			expected: "/path/to/repo",
			wantErr:  false,
		},
		{
			name: "failure",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"rev-parse", "--show-toplevel"}},
					Out:    "",
					Err:    errors.New("not a git repository"),
				},
			},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//
			mockExe := testutils.NewMockExecutor(tt.mocks)
			//new(testutils.MockExecutor)
			//for _, mock := range tt.mocks {
			//		mockExe.On(mock.Method, mock.Args.([]interface{})...).Return(mock.Out, mock.Err)
			//	}

			got, err := utils.GetGitRootDirectory(mockExe)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}

			mockExe.AssertExpectations(t)
		})
	}
}

func TestSetWorkDirToGitHubRoot(t *testing.T) {
	tests := []struct {
		name    string
		mocks   []testutils.MockContent
		wantErr bool
	}{
		{
			name: "success",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"remote", "get-url", "origin"}},
					Out:    "git@github.com:example/test.git",
					Err:    nil,
				},
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"rev-parse", "--show-toplevel"}},
					Out:    "/path/to/repo\n",
					Err:    nil,
				},
				{
					Method: "Chdir",
					Args:   []interface{}{"/path/to/repo"},
					Out:    "",
					Err:    nil,
				},
			},
			wantErr: false,
		},
		{
			name: "failure: git remote returns error",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"remote", "get-url", "origin"}},
					Out:    "",
					Err:    errors.New("not a git repository"),
				},
			},
			wantErr: true,
		},
		{
			name: "failure: git remote returns invalid url",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"git", []string{"remote", "get-url", "origin"}},
					Out:    "thiswasnotagithuburl",
					Err:    nil,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExe := new(testutils.MockExecutor)
			for _, mock := range tt.mocks {
				mockExe.On(mock.Method, mock.Args.([]interface{})...).Return(mock.Out, mock.Err)
			}

			err := utils.SetWorkDirToGitHubRoot(mockExe)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			mockExe.AssertExpectations(t)
		})
	}
}

func TestListFilesInDirectory(t *testing.T) {
	tests := []struct {
		name     string
		mocks    []testutils.MockContent
		expected []string
		wantErr  bool
	}{
		{
			name: "success",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"ls", []string{"/path/to/dir"}},
					Out:    "file1\nfile2\n",
					Err:    nil,
				},
			},
			expected: []string{"file1", "file2"},
			wantErr:  false,
		},
		{
			name: "failure",
			mocks: []testutils.MockContent{
				{
					Method: "Command",
					Args:   []interface{}{"ls", []string{"/path/to/dir"}},
					Out:    "",
					Err:    errors.New("directory does not exist"),
				},
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExe := new(testutils.MockExecutor)
			for _, mock := range tt.mocks {
				mockExe.On(mock.Method, mock.Args.([]interface{})...).Return(mock.Out, mock.Err)
			}

			got, err := utils.ListFilesInDirectory(mockExe, "/path/to/dir")
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}

			mockExe.AssertExpectations(t)
		})
	}
}
