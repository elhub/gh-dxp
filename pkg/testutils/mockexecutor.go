package testutils

import (
	"bytes"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockExecutor is a mock implementation of the Executor interface.
type MockExecutor struct {
	mock.Mock
}

// Command pretends to run an OS command and returns its output.
func (m *MockExecutor) Command(name string, arg ...string) (string, error) {
	args := m.Called(name, arg)
	return args.String(0), args.Error(1)
}

// CommandContext pretends to run an OS command with a context and returns an error.
func (m *MockExecutor) CommandContext(ctx context.Context, name string, arg ...string) error {
	args := m.Called(ctx, name, arg)
	return args.Error(1)
}

// GH pretends to run a GitHub CLI command and returns its output.
func (m *MockExecutor) GH(arg ...string) (bytes.Buffer, error) {
	args := m.Called(arg)
	return *bytes.NewBufferString(args.String(0)), args.Error(1)
}

// Chdir pretends to change the current working directory.
func (m *MockExecutor) Chdir(dir string) error {
	args := m.Called(dir)
	return args.Error(1)
}

// MockContent represents the content of a mock method call to MockExecutor.
type MockContent struct {
	Method string
	Args   interface{}
	Out    string
	Err    error
}
