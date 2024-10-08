package test_test

import (
	"errors"
	"testing"

	"github.com/elhub/gh-dxp/pkg/test"
	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name           string
		gitRoot        string
		gitRootError   error
		expectedResult bool
		expectedErr    error
		testFile       string
	}{
		{
			name:           "Test Makefile",
			gitRoot:        "/home/repo-name",
			testFile:       "/home/repo-name/Makefile",
			expectedResult: true,
		},
		{
			name:           "Test Gradlew",
			gitRoot:        "/home/repo-name",
			testFile:       "/home/repo-name/gradlew",
			expectedResult: true,
		},
		{
			name:           "Test npm",
			gitRoot:        "/home/repo-name",
			testFile:       "/home/repo-name/package.json",
			expectedResult: true,
		},
		{
			name:           "Test maven",
			gitRoot:        "/home/repo-name",
			testFile:       "/home/repo-name/pom.xml",
			expectedResult: true,
		},
		{
			name:           "Failing test",
			gitRoot:        "/home/repo-name",
			testFile:       "/home/repo-name/Makefile",
			expectedResult: false,
			expectedErr:    errors.New("failed tests"),
		},
		{
			name:           "No test file",
			gitRoot:        "/home/repo-name",
			expectedResult: false,
		},
		{
			name:           "Not in git repo",
			gitRootError:   errors.New("Not a git repo"),
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.FileExists = func(path string) bool {
				return path == tt.testFile
			}

			mockExe := new(testutils.MockExecutor)

			mockExe.On("Command", "git", []string{"rev-parse", "--show-toplevel"}).Return(tt.gitRoot, tt.gitRootError)
			mockExe.On("CommandContext", mock.Anything, "./gradlew", []string{"test"}).Return(nil, tt.expectedErr)
			mockExe.On("CommandContext", mock.Anything, "make", []string{"check"}).Return(nil, tt.expectedErr)
			mockExe.On("CommandContext", mock.Anything, "npm", []string{"test"}).Return(nil, tt.expectedErr)
			mockExe.On("CommandContext", mock.Anything, "mvn", []string{"test"}).Return(nil, tt.expectedErr)

			res, err := test.RunTest(mockExe)

			assert.Equal(t, tt.expectedResult, res)
			if tt.expectedErr != nil || tt.gitRootError != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
