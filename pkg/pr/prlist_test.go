package pr_test

import (
	"testing"

	"github.com/elhub/gh-dxp/pkg/pr"
	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecuteList(t *testing.T) {

	prDetails := "{\"additions\":1,\"author\":{\"id\":\"U_kgDOCXKgZ\",\"is_bot\":false," +
		"\"login\":\"my-user-statnett\",\"name\":\"Stat Nett\"},\"createdAt\":" +
		"\"2024-08-13T08:23:12Z\",\"deletions\":1,\"headRepository\":{\"id\":\"R_kgDOMCGoz\",\"name\"" +
		":\"gh-xyz\"},\"number\":1,\"reviewDecision\":\"\",\"title\":\"the title of this here pr\"}"

	tests := []struct {
		name                   string
		expectedErr            error
		authorSearchResponse   string
		prDetailResponse       string
		reviewerSearchResponse string
	}{
		{
			name:                   "User is author of a PR",
			authorSearchResponse:   "[{\"number\":1,\"repository\":{\"name\":\"gh-xyz\",\"nameWithOwner\":\"elhub/gh-xyz\"}}]",
			reviewerSearchResponse: "[]",
			prDetailResponse:       prDetails,
		},
		{
			name:                   "User has one review requested",
			authorSearchResponse:   "[]",
			reviewerSearchResponse: "[{\"number\":1,\"repository\":{\"name\":\"gh-xyz\",\"nameWithOwner\":\"elhub/gh-xyz\"}}]",
			prDetailResponse:       prDetails,
		},
		{
			name:                   "There are no PR's assigned to the user",
			reviewerSearchResponse: "[]",
			authorSearchResponse:   "[]",
		},
		{
			name:                 "Author search fails",
			authorSearchResponse: "[]",
			expectedErr:          errors.New("something went wrong during author search"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExe := new(testutils.MockExecutor)

			mockExe.On("GH", []string{"search", "prs", "--author=@me", "--state=open", "--json", "number,repository"}).Return(tt.authorSearchResponse, tt.expectedErr)
			mockExe.On("GH", []string{"search", "prs", "--review-requested=@me", "--state=open", "--json", "number,repository"}).Return(tt.reviewerSearchResponse, tt.expectedErr)

			mockExe.On("GH", []string{"pr", "view", "https://github.com/elhub/gh-xyz/pull/1", "--json", "additions,author,createdAt,deletions,headRepository,number,title,reviewDecision"}).Return(tt.prDetailResponse, tt.expectedErr)

			opts := &pr.ListOptions{TestRun: true, Mine: true, ReviewRequested: true}

			err := pr.ExecuteList(mockExe, opts)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.True(t, errors.As(err, &tt.expectedErr))
			} else {
				require.NoError(t, err)
			}
		})
	}
}
