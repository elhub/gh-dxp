package merge_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	merge "github.com/elhub/gh-dxp/pkg/prmerge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockExecutor struct {
	mock.Mock
}

func (m *MockExecutor) Command(name string, arg ...string) (string, error) {
	args := m.Called(name, arg)
	return args.String(0), args.Error(1)
}

func (m *MockExecutor) CommandContext(ctx context.Context, name string, arg ...string) error {
	args := m.Called(ctx, name, arg)
	return args.Error(1)
}

func (m *MockExecutor) GH(arg ...string) (bytes.Buffer, error) {
	args := m.Called(arg)
	return *bytes.NewBufferString(args.String(0)), args.Error(1)
}

func TestExecute(t *testing.T) {
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
			mockExe := new(MockExecutor)
			mockExe.On("Command", "git", []string{"branch", "--show-current"}).Return(tt.pushBranch, tt.pushBranchErr)
			mockExe.On("GH", []string{"pr", "list", "-H", tt.pushBranch, "--json", "number", "--jq", ".[].number"}).
				Return(tt.prNumber, tt.prNumberErr)
			mockExe.On("GH", []string{"pr", "view", "--json", "title", "--jq", ".title"}).Return(tt.prTitle, tt.prTitleErr)
			mockExe.On("GH", []string{"pr", "merge", "--squash", "--delete-branch"}).Return(tt.prMerge, tt.prMergeErr)

			err := merge.Execute(mockExe, &merge.Options{
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
