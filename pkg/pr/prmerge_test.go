package pr_test

import (
	"errors"
	"testing"

	"github.com/elhub/gh-dxp/pkg/pr"
	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecuteMerge(t *testing.T) {
	tests := []struct {
		name          string
		pushBranch    string
		pushBranchErr error
		prNumber      string
		prNumberErr   error
		prTitle       string
		prTitleErr    error
		prMerge       string
		prMergeErr    error
		expectedErr   error
	}{
		{
			name:        "Test successful PR creation",
			pushBranch:  "branch1",
			prNumber:    "3",
			prTitle:     "PR title",
			prMerge:     "pull request merged",
			expectedErr: nil,
		},
		{
			name:          "Test error getting current branch",
			pushBranch:    "",
			pushBranchErr: errors.New("error getting current branch"),
			expectedErr:   errors.New("error getting current branch"),
		},
		{
			name:        "Test error getting PR number",
			pushBranch:  "branch1",
			prNumberErr: errors.New("error getting PR number"),
			expectedErr: errors.New("Failed to find existing PR"),
		},
		{
			name:        "Test error getting PR title",
			pushBranch:  "branch1",
			prNumber:    "3",
			prTitleErr:  errors.New("Error getting PR title"),
			expectedErr: errors.New("Error getting PR title"),
		},
		{
			name:        "Test error merging PR",
			pushBranch:  "branch1",
			prNumber:    "3",
			prTitle:     "PR title",
			prMergeErr:  errors.New("error merging PR"),
			expectedErr: errors.New("Failed to merge pull request #3"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExe := new(testutils.MockExecutor)
			mockExe.On("Command", "git", []string{"branch", "--show-current"}).Return(tt.pushBranch, tt.pushBranchErr)
			mockExe.On("GH", []string{"pr", "list", "-H", tt.pushBranch, "--json", "number", "--jq", ".[].number"}).
				Return(tt.prNumber, tt.prNumberErr)
			mockExe.On("GH", []string{"pr", "view", "--json", "title", "--jq", ".title"}).Return(tt.prTitle, tt.prTitleErr)
			mockExe.On("GH", []string{"pr", "merge", "--squash", "--delete-branch"}).Return(tt.prMerge, tt.prMergeErr)

			err := pr.ExecuteMerge(mockExe, &pr.MergeOptions{
				AutoConfirm: true,
			})

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
