package pr

import (
	"testing"

	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestInferLabelFromCommits(t *testing.T) {
	tests := []struct {
		name         string
		commitOutput string
		commitErr    error
		expected     string
	}{
		{
			name:         "maps feat to Feature",
			commitOutput: "feat: add new endpoint",
			expected:     "Feature",
		},
		{
			name:         "maps scoped fix to Bugfix",
			commitOutput: "fix(api): handle nil response",
			expected:     "Bugfix",
		},
		{
			name:         "maps chore to Chore",
			commitOutput: "chore: bump dependency",
			expected:     "Chore",
		},
		{
			name:         "maps docs to Documentation",
			commitOutput: "docs: update guide",
			expected:     "Documentation",
		},
		{
			name:         "maps refactor to Refactor",
			commitOutput: "refactor: simplify parser",
			expected:     "Refactor",
		},
		{
			name:         "maps style to Style",
			commitOutput: "style: reformat handlers",
			expected:     "Style",
		},
		{
			name:         "maps test to Test",
			commitOutput: "test: add parser coverage",
			expected:     "Test",
		},
		{
			name:         "maps build to Build",
			commitOutput: "build: update ci image",
			expected:     "Build",
		},
		{
			name:         "returns empty for unknown conventional type",
			commitOutput: "ci: update pipeline",
			expected:     "",
		},
		{
			name:         "returns empty for non conventional subject",
			commitOutput: "update readme",
			expected:     "",
		},
		{
			name:         "returns empty on git log error",
			commitErr:    errors.New("git log failed"),
			expected:     "",
		},
		{
			name:         "uses only first commit line",
			commitOutput: "fix: repair parser\nfeat: add parser option",
			expected:     "Bugfix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExe := new(testutils.MockExecutor)
			mockExe.On("Command", "git", []string{"log", "main..feature-branch", "--oneline", "--pretty=format:%s"}).
				Return(tt.commitOutput, tt.commitErr)

			pr := PullRequest{
				branchID:     "feature-branch",
				targetBranch: "main",
			}

			got := inferLabelFromCommits(mockExe, pr)
			assert.Equal(t, tt.expected, got)
			mockExe.AssertExpectations(t)
		})
	}
}
