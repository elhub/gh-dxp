package repo_test

import (
	"errors"
	"testing"
	"time"

	"github.com/elhub/gh-dxp/pkg/repo"
	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func TestRun_ExecuteClone(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		options *repo.Options
	}{
		{
			name:    "Clone all without dryrun",
			pattern: "",
			options: &repo.Options{DryRun: false},
		},
		{
			name:    "Clone all with dryrun",
			pattern: "",
			options: &repo.Options{DryRun: true},
		},
		{
			name:    "Clone with pattern without dryrun",
			pattern: "repo",
			options: &repo.Options{DryRun: false},
		},
		{
			name:    "Clone with pattern with dryrun",
			pattern: "repo",
			options: &repo.Options{DryRun: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExe := new(testutils.MockExecutor)
			mockExe.On("GH", []string{"api", "user/orgs"}).Return(`[{"login": "myorg","id": 3679327}]`, nil)
			if tt.pattern == "" {
				mockExe.On("GH", []string{"search", "repos", "--archived=false", "--json", "name,fullName,url", "--limit=1000", "--owner", "myorg"}).Return(`[{"name": "repo1", "fullName": "myorg/repo1", "url": "https://github.com/myorg/repo1"}, {"name": "repo2", "fullName": "myorg/repo2", "url": "https://github.com/myorg/repo2"}]`, nil)
			} else {
				mockExe.On("GH", []string{"search", "repos", tt.pattern, "--match", "name", "--archived=false", "--json", "name,fullName,url", "--limit=1000", "--owner", "myorg"}).Return(`[{"name": "repo1", "fullName": "myorg/repo1", "url": "https://github.com/myorg/repo1"}, {"name": "repo2", "fullName": "myorg/repo2", "url": "https://github.com/myorg/repo2"}]`, nil)
			}
			if !tt.options.DryRun { // If not dry run, we expect the repos to be cloned
				mockExe.On("GH", []string{"repo", "clone", "myorg/repo1"}).Return("", errors.New("Mocked error")).Once() //This is intended to test the retry logic
				mockExe.On("GH", []string{"repo", "clone", "myorg/repo1"}).Return("", nil)
				mockExe.On("GH", []string{"repo", "clone", "myorg/repo2"}).Return("", nil)
			}

			err := repo.ExecuteClone(mockExe, tt.pattern, mockSleep, tt.options)
			assert.NoError(t, err)
			mockExe.AssertExpectations(t)
		})
	}
}

func mockSleep(_ time.Duration) {}
