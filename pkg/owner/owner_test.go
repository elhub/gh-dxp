package owner_test

import (
	"testing"

	"github.com/elhub/gh-dxp/pkg/owner"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/stretchr/testify/require"
)

/*
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
*/

func TestExecute(t *testing.T) {
	tests := []struct {
		path   string
		owners []string
	}{
		{
			path:   "README.md",
			owners: []string{"@elhub/devxp"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			/*mockExe := new(MockExecutor)

			mockExe.On("Command", "git", []string{"rev-parse", "--show-toplevel"}).Return(tt.gitRoot, tt.gitRootError)
			mockExe.On("CommandContext", mock.Anything, "gradlew", []string{"test"}).Return(nil, tt.expectedErr)
			mockExe.On("CommandContext", mock.Anything, "make", []string{"check"}).Return(nil, tt.expectedErr)
			mockExe.On("CommandContext", mock.Anything, "npm", []string{"test"}).Return(nil, tt.expectedErr)
			mockExe.On("CommandContext", mock.Anything, "mvn", []string{"test"}).Return(nil, tt.expectedErr)*/

			exe := utils.LinuxExecutor()

			owners, err := owner.Execute(tt.path, exe)
			require.NoError(t, err)
			require.Equal(t, tt.owners, owners)
		})
	}
}
