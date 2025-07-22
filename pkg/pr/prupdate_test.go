package pr_test

import (
	"testing"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/pr"
	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestExecuteUpdate(t *testing.T) {
	tests := []struct {
		name             string
		currentBranch    string
		currentBranchErr error
		pushBranch       string
		pushBranchErr    error
		prListNumber     string
		prListNErr       error
		prListURL        string
		prListUErr       error
		gitLog           string
		gitLogErr        error
		repoBranchName   string
		repoBranchErr    error
		prCreate         string
		prCreateErr      error
		expectedErr      error
		currentChanges   string
		expectedLintErr  error
		modifiedFiles    string
		existingBranches string
	}{
		{
			name:             "Test successful PR update",
			currentBranch:    "branch1",
			pushBranch:       "branch1",
			prListNumber:     "3",
			expectedErr:      nil,
			modifiedFiles:    "pkg/cmd/lint.go\npkg/lint/lint.go\n",
			existingBranches: "main\ndifferentBranch\n",
			currentChanges:   "M  pkg/cmd/lint.go\nM  pkg/lint/lint.go\n",
		},
		{
			name:             "Test error in getting current branch",
			currentBranchErr: errors.New("error getting current branch"),
			expectedErr:      errors.New("error getting current branch"),
			modifiedFiles:    "pkg/cmd/lint.go\npkg/lint/lint.go\n",
			existingBranches: "main\ndifferentBranch\n",
		},
		{
			name:             "Test error in checking for existing PR",
			currentBranch:    "branch1",
			pushBranch:       "branch1",
			prListNumber:     "",
			prListNErr:       errors.New("error checking for existing PR"),
			expectedErr:      errors.New("Failed to find existing PR"),
			existingBranches: "main\ndifferentBranch\n",
			modifiedFiles:    "pkg/cmd/lint.go\npkg/lint/lint.go\n",
			currentChanges:   "M  pkg/cmd/lint.go\nM  pkg/lint/lint.go\n",
		},
		{
			name:             "Test error in update flow - git push",
			currentBranch:    "branch1",
			pushBranch:       "branch1",
			pushBranchErr:    errors.New("error pushing branch"),
			prListNumber:     "1",
			prListURL:        "https://github.com/elhub/demo/pull/3",
			expectedErr:      errors.New("error pushing branch"),
			existingBranches: "main\ndifferentBranch\n",
			modifiedFiles:    "pkg/cmd/lint.go\npkg/lint/lint.go\n",
			currentChanges:   "M  pkg/cmd/lint.go\nM  pkg/lint/lint.go\n",
		},
		{
			name:             "Test error in update flow - list URL",
			currentBranch:    "branch1",
			pushBranch:       "branch1",
			prListNumber:     "1",
			prListNErr:       nil,
			prListURL:        "",
			prListUErr:       errors.New("error getting PR URL"),
			expectedErr:      errors.New("error getting PR URL"),
			existingBranches: "main\ndifferentBranch\n",
			modifiedFiles:    "pkg/cmd/lint.go\npkg/lint/lint.go\n",
			currentChanges:   "M  pkg/cmd/lint.go\nM  pkg/lint/lint.go\n",
		},
		{
			name:             "Test local has untracked changes",
			currentBranch:    "branch1",
			pushBranch:       "branch1",
			prListNumber:     "1",
			prListNErr:       nil,
			prListURL:        "https://github.com/elhub/demo/pull/3",
			gitLog:           "commit 1",
			repoBranchName:   "main",
			prCreate:         "pull request created",
			expectedErr:      errors.New("No tracked changes found, skipping commit"),
			currentChanges:   "?? untracked_change.go",
			existingBranches: "main\ndifferentBranch\n",
			modifiedFiles:    "pkg/cmd/lint.go\npkg/lint/lint.go\n",
		},
		{
			name:             "Test local has tracked changes",
			currentBranch:    "branch1",
			pushBranch:       "branch1",
			prListNumber:     "1",
			prListNErr:       nil,
			prListURL:        "https://github.com/elhub/demo/pull/3",
			gitLog:           "commit 1",
			repoBranchName:   "main",
			prCreate:         "pull request created",
			expectedErr:      nil,
			currentChanges:   "M  tracked_change.go",
			existingBranches: "main\ndifferentBranch\n",
			modifiedFiles:    "pkg/cmd/lint.go\npkg/lint/lint.go\n",
		},
		{
			name:             "Test local has multiple tracked changes",
			currentBranch:    "branch1",
			pushBranch:       "branch1",
			prListNumber:     "1",
			prListNErr:       nil,
			prListURL:        "https://github.com/elhub/demo/pull/3",
			gitLog:           "commit 1",
			repoBranchName:   "main",
			prCreate:         "pull request created",
			expectedErr:      nil,
			existingBranches: "main\ndifferentBranch\n",
			currentChanges:   " M tracked_change.go\n M tracked_change2.go",
		},
		{
			name:             "Test lint is failing",
			expectedLintErr:  errors.New("exit status 1"),
			expectedErr:      errors.New("exit status 1"),
			modifiedFiles:    "pkg/cmd/lint.go\npkg/lint/lint.go\n",
			existingBranches: "main\ndifferentBranch\n",
			currentBranch:    "branch1",
			prListNumber:     "1",
			currentChanges:   "M  pkg/cmd/lint.go\nM  pkg/lint/lint.go\n",
		},
		{
			name:             "Test lint is failing",
			expectedErr:      errors.New("No PR found for branch branch1"),
			modifiedFiles:    "pkg/cmd/lint.go\npkg/lint/lint.go\n",
			existingBranches: "main\ndifferentBranch\n",
			currentBranch:    "branch1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			linterArgs := []string{"mega-linter-runner", "--flavor", "cupcake",
				"-e", "LINTER_RULES_PATH=/tmp",
				"-e", "MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml"}

			linterArgs = append(linterArgs, "--filesonly")
			linterArgs = append(linterArgs, utils.ConvertTerminalOutputIntoList(tt.modifiedFiles)...)

			mockExe := new(testutils.MockExecutor)
			mockExe.On("Command", "git", []string{"rev-parse", "--show-toplevel"}).Return("/home/repo-name", nil)
			mockExe.On("Command", "git", []string{"status", "--porcelain"}).Return(tt.currentChanges, nil)
			mockExe.On("Command", "git", []string{"branch", "--show-current"}).Return(tt.currentBranch, tt.currentBranchErr)
			mockExe.On("Command", "git", []string{"remote", "set-head", "origin", "--auto"}).Return("", nil)
			mockExe.On("Command", "git", []string{"symbolic-ref", "--short", "refs/remotes/origin/HEAD"}).Return("origin/main", nil)
			mockExe.On("Command", "git", []string{"diff", "--name-only", "origin/main", "--relative"}).Return(tt.modifiedFiles, nil)
			mockExe.On("Command", "git", []string{"push"}).Return(tt.pushBranch, tt.pushBranchErr)
			mockExe.On("Command", "git", []string{"branch"}).Return(tt.existingBranches, nil)
			mockExe.On("Command", "git", []string{"add", "-u"}).Return("", nil)
			mockExe.On("Command", "git", []string{"fetch", "origin", "main"}).Return("", nil)
			mockExe.On("Command", "git", []string{"commit", "-m", "default commit message"}).Return("", nil)
			mockExe.On("Command", "git", []string{"show-ref", "--verify", "--quiet", "refs/heads/" + tt.currentBranch})
			mockExe.On("CommandContext", mock.Anything, "npx", linterArgs).Return("", tt.expectedLintErr)
			mockExe.On("Command", "git", []string{"push", "--set-upstream", "origin", tt.currentBranch}).
				Return(tt.pushBranch, tt.pushBranchErr)
			mockExe.On("Command", "git", []string{"log", "main.." + tt.currentBranch, "--oneline", "--pretty=format:%s"}).
				Return(tt.gitLog, tt.gitLogErr)
			mockExe.On("GH", []string{"pr", "list", "-H", tt.currentBranch, "--json", "number", "--jq", ".[].number"}).
				Return(tt.prListNumber, tt.prListNErr)
			mockExe.On("GH", []string{"pr", "list", "-H", tt.currentBranch, "--json", "url", "--jq", ".[].url"}).
				Return(tt.prListURL, tt.prListUErr)
			mockExe.On("GH", []string{"pr", "create", "--title", tt.gitLog, "--body", "Testing:\n- [ ] Unit Tests\n" +
				"- [ ] Integration Tests\n\n\nDocumentation:\n- No updates\n", "--base", "main"}).
				Return(tt.prCreate, tt.prCreateErr)
			mockExe.On("GH", []string{"repo", "view", "--json", "defaultBranchRef", "--jq", ".defaultBranchRef.name"}).
				Return(tt.repoBranchName, tt.repoBranchErr)

			err := pr.ExecuteUpdate(mockExe,
				&config.Settings{},
				&pr.UpdateOptions{
					TestRun: true,
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
