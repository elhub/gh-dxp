package renovate_test

import (
	"testing"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/renovate"
	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const MainBranch = "origin/main"


func TestRenovateValidationError(t *testing.T) {
	test := []struct {
		name	  		string
		filesChanged	string
		forceValidate	bool
		expectedErr 	error
		expectedRun		bool
	}{
		{
			name:	 		"Test run on renovate config change",
			filesChanged:	".github/renovate.json",
			forceValidate: 	false,
			expectedErr:	nil,
			expectedRun: 	true,
		},
		{
			name:	  		"Test don't run with no renovate config changes",
			filesChanged:	"file1.txt\nfile2.txt",
			forceValidate: 	false,
			expectedErr:	nil,
			expectedRun: 	false,
		},
		{
			name:			"Test run with no renovate config changes and --force used",
			filesChanged:	"file1.txt\nfile2.txt",
			forceValidate: 	true,
			expectedErr:	nil,
			expectedRun: 	true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			mockExe := new(testutils.MockExecutor)
			mockExe.On("Command", "git", []string{"branch"}).Return("main\ndifferentBranch\n", nil)
			mockExe.On("Command", "git", []string{"fetch", "origin", "main"}).Return("", nil)
			mockExe.On("Command", "git", []string{"remote", "set-head", "origin", "--auto"}).Return("", nil)
			mockExe.On("Command", "git", []string{"symbolic-ref", "--short", "refs/remotes/origin/HEAD"}).Return("origin/main", nil)
			mockExe.On("Command", "git", []string{"diff", "--name-only", MainBranch, "--relative"}).Return(tt.filesChanged, nil)

			renovateArgs := []string{
				"--package", "renovate@43.78.0", "renovate-config-validator", "--strict",
			}

			if (tt.expectedRun) {
				mockExe.On("CommandContext", mock.Anything, "npx", renovateArgs).Return(nil, nil)
			}

			err := renovate.Run(mockExe, &config.Settings{}, &renovate.Options{Force: tt.forceValidate})

			if tt.expectedErr != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if (!tt.expectedRun) {
				mockExe.AssertNotCalled(t, "CommandContext", mock.Anything, "npx", renovateArgs)
			}

			mockExe.AssertExpectations(t)
		})
	}
}
