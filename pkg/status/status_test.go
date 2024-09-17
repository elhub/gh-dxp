package status_test

import (
	"testing"

	"github.com/elhub/gh-dxp/pkg/status"
	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusAll(t *testing.T) {
	tests := []struct {
		name        string
		expected    string
		expectedErr error
		all         bool
		repo        bool
		pr          bool
		branches    bool
		issue       bool
	}{
		{
			name:        "Test PRs",
			expected:    "PR Status:\nRelevant pull requests in repo-name\n",
			expectedErr: nil,
			pr:          true,
		},
		{
			name:        "Test branches",
			expected:    "branch1\nbranch2\n",
			expectedErr: nil,
			branches:    true,
		},
		{
			name:        "Test repo",
			expected:    "Repository: git@github.com:elhub/repo-name.git\n",
			expectedErr: nil,
			repo:        true,
		},
		{
			name:        "Test issue",
			expected:    "Assigned PRs/Review Requests:\n",
			expectedErr: nil,
			issue:       true,
		},
		{
			name:        "Test all",
			expected:    "Repository: git@github.com:elhub/repo-name.git\nPR Status:",
			expectedErr: nil,
			all:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExec := new(testutils.MockExecutor)
			mockExec.On("Command", "git", []string{"remote", "get-url", "origin"}).
				Return("git@github.com:elhub/demo.git", nil)
			mockExec.On("GH", []string{"pr", "status"}).Return(tt.expected, nil)
			mockExec.On("Command", "git", []string{"branch", "-a"}).Return(tt.expected, nil)
			mockExec.On("GH", []string{"issue", "status"}).Return(tt.expected, nil)

			err := status.Execute(mockExec,
				&status.Options{
					All:      tt.all,
					Repo:     tt.repo,
					Pr:       tt.pr,
					Branches: tt.branches,
					Issue:    tt.issue,
				},
			)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
