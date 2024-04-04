package diff_test

import (
	"bytes"
	"testing"

	"github.com/elhub/gh-dxp/pkg/diff"
	"github.com/pkg/errors"
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

func (m *MockExecutor) GH(arg ...string) (bytes.Buffer, bytes.Buffer, error) {
	args := m.Called(arg)
	return *bytes.NewBufferString(args.String(0)), bytes.Buffer{}, args.Error(1)
}

/*
func TestExecute(t *testing.T) {
	tests := []struct {
		name             string
		currentBranch    string
		currentBranchErr error
		pushBranch       string
		pushBranchErr    error
		prListNumber     string
		prListNErr       error
		prListUrl        string
		prListUErr       error
		gitLog           string
		gitLogErr        error
		repoBranchName   string
		repoBranchErr    error
		prCreate         string
		prCreateErr      error
		expectedErr      error
	}{
		/*
			{
				name:           "Test successful PR creation",
				currentBranch:  "branch1",
				pushBranch:     "branch1",
				prListNErr:     errors.New("no pull requests match your search in elhub/demo"),
				prListUrl:      "https://github.com/elhub/demo/pull/3",
				gitLog:         "commit 1",
				repoBranchName: "main",
				prCreate:       "pull request created",
				expectedErr:    nil,
			},
			/*	{
						name:          "Test successful PR update",
						currentBranch: "branch1",
						pushBranch:    "branch1",
						prListNumber:  "3",
						expectedErr:   nil,
					},
					{
						name:             "Test error in getting current branch",
						currentBranchErr: errors.New("error getting current branch"),
						expectedErr:      errors.New("error getting current branch"),
					},
		{
			name:             "Test error in checking for existing PR",
			currentBranch:    "branch1",
			currentBranchErr: nil,
			pushBranch:       "branch1",
			pushBranchErr:    nil,
			prListNumber:     "",
			prListNErr:       errors.New("error checking for existing PR"),
			prListUrl:        "",
			prListUErr:       nil,
			gitLog:           "commit 1",
			gitLogErr:        nil,
			repoBranchName:   "main",
			repoBranchErr:    nil,
			prCreate:         "",
			prCreateErr:      nil,
			expectedErr:      errors.New("error checking for existing PR"),
		},

		/*
			{
				name:             "Test error in checking for existing PR",
				currentBranch:    "branch1",
				currentBranchErr: nil,
				pushBranch:       "branch1",
				pushBranchErr:    nil,
				prListNumber:     "",
				prListNErr:       errors.New("Failed to find existing PR"),
				prListUrl:        "",
				prListUErr:       nil,
				gitLog:           "commit 1",
				gitLogErr:        nil,
				repoBranchName:   "main",
				repoBranchErr:    nil,
				prCreate:         "",
				prCreateErr:      nil,
				expectedErr:      errors.New("Failed to find existing PR"),
			},
			/*
				{
					name:          "Test error in updating PR",
					currentBranch: "branch1",
					prId:          "pr1",
					updateErr:     errors.New("error updating PR"),
					expectedErr:   errors.New("error updating PR"),
				},
				{
					name:          "Test error in creating PR",
					currentBranch: "branch1",
					prId:          "",
					createErr:     errors.New("error creating PR"),
					expectedErr:   errors.New("error creating PR"),
				},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExe := new(MockExecutor)
			mockExe.On("Command", "git", []string{"branch", "--show-current"}).Return(tt.currentBranch, tt.currentBranchErr)
			mockExe.On("Command", "git", []string{"push"}).Return(tt.pushBranch, tt.pushBranchErr)
			mockExe.On("Command", "git", []string{"push", "--set-upstream", "origin", tt.currentBranch}).Return(tt.pushBranch, tt.pushBranchErr)
			mockExe.On("Command", "git", []string{"log", "main.." + tt.currentBranch, "--oneline", "--pretty=format:%s"}).Return(tt.gitLog, tt.gitLogErr)
			mockExe.On("GH", []string{"pr", "list", "-H", tt.currentBranch, "--json", "number", "--jq", ".[].number"}).Return(tt.prListNumber, nil, errors.New("Hello"))
			mockExe.On("GH", []string{"pr", "list", "-H", tt.currentBranch, "--json", "url", "--jq", ".[].url"}).Return(tt.prListUrl, nil, tt.prListUErr)
			mockExe.On("GH", []string{"pr", "create", "--title", tt.gitLog, "--body", "Testing:\n- [ ] Unit Tests\n- [ ] Integration Tests\n- Test Command: \n\nDocumentation:\n- No updates\n", "--base", "main"}).Return(tt.prCreate, nil, tt.prCreateErr)
			mockExe.On("GH", []string{"repo", "view", "--json", "defaultBranchRef", "--jq", ".defaultBranchRef.name"}).Return(tt.repoBranchName, nil, tt.repoBranchErr)

			err := diff.Execute(mockExe, nil, &diff.Options{
				AutoConfirm: true,
			})

			mockExe.AssertCalled(t, "GH", []string{"pr", "list", "-H", "branch1", "--json", "number", "--jq", ".[].number"})

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}*/

func TestCheckForExistingPR(t *testing.T) {
	tests := []struct {
		name        string
		branchId    string
		expectedNum string
		expectedErr error
	}{
		{
			name:        "Test existing PR",
			branchId:    "branch1",
			expectedNum: "3",
			expectedErr: nil,
		},
		{
			name:        "Test no existing PR",
			branchId:    "branch2",
			expectedNum: "",
			expectedErr: errors.New("Failed to find existing PR"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExe := new(MockExecutor)
			mockExe.On("GH", []string{"pr", "list", "-H", tt.branchId, "--json", "number", "--jq", ".[].number"}).Return(tt.expectedNum, nil, tt.expectedErr)

			num, err := diff.CheckForExistingPR(mockExe, tt.branchId)

			mockExe.AssertCalled(t, "GH", []string{"pr", "list", "-H", tt.branchId, "--json", "number", "--jq", ".[].number"})

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedNum, num)
		})
	}
}
