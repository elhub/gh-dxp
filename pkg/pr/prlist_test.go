package pr_test

import (
	"bytes"
	"testing"

	"github.com/elhub/gh-dxp/pkg/pr"
	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecuteList(t *testing.T) {
	tests := []struct {
		name           string
		expectedErr    error
		searchResponse string
	}{
		{
			name:           "",
			expectedErr:    nil,
			searchResponse: ("[{'number':1,'repository':{'name':'gh-xyz','nameWithOwner':'elhub/gh-xyz}}]"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExe := new(testutils.MockExecutor)

			var responseBuffer bytes.Buffer
			responseBuffer.WriteString(tt.searchResponse)

			//mockExe.On("Command", "git", []string{"branch", "--show-current"}).Return(tt.pushBranch, tt.pushBranchErr)
			mockExe.On("GH", []string{"search", "prs", "--author=@me", "--state=open", "--json", "number,repository"}).Return(responseBuffer, tt.expectedErr)

			opts := &pr.ListOptions{TestRun: true, Mine: true, ReviewRequested: true}

			err := pr.ExecuteList(mockExe, opts)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
