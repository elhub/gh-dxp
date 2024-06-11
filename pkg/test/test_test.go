package test_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/elhub/gh-dxp/pkg/test"
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
		name         string
		gitRoot      string
		gitRootError error
		expectedErr  error
		testFile     string
	}{
		{
			name:     "Test Makefile",
			gitRoot:  "/home/repo-name",
			testFile: "/home/repo-name/Makefile",
		},
		{
			name:     "Test Gradlew",
			gitRoot:  "/home/repo-name",
			testFile: "/home/repo-name/Gradlew",
		},
		{
			name:     "Test npm",
			gitRoot:  "/home/repo-name",
			testFile: "/home/repo-name/package.json",
		},
		{
			name:     "Test maven",
			gitRoot:  "/home/repo-name",
			testFile: "/home/repo-name/pom.xml",
		},
		{
			name:        "Failing test",
			gitRoot:     "/home/repo-name",
			testFile:    "/home/repo-name/Makefile",
			expectedErr: errors.New("failed tests"),
		},
		{
			name:    "No test file",
			gitRoot: "/home/repo-name",
		},
		{
			name:         "Not in git repo",
			gitRootError: errors.New("Not a git repo"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			test.FileExists = func(path string) bool {
				return path == tt.testFile
			}

			mockExe := new(MockExecutor)

			mockExe.On("Command", "git", []string{"rev-parse", "--show-toplevel"}).Return(tt.gitRoot, tt.gitRootError)
			mockExe.On("CommandContext", mock.Anything, "gradlew", []string{"test"}).Return(nil, tt.expectedErr)
			mockExe.On("CommandContext", mock.Anything, "make", []string{"test"}).Return(nil, tt.expectedErr)
			mockExe.On("CommandContext", mock.Anything, "npm", []string{"test"}).Return(nil, tt.expectedErr)
			mockExe.On("CommandContext", mock.Anything, "mvn", []string{"test"}).Return(nil, tt.expectedErr)

			err := test.RunTest(mockExe)

			if tt.expectedErr != nil || tt.gitRootError != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
