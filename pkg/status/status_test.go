package status_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/elhub/gh-dxp/pkg/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockExecutor struct {
    mock.Mock
}

func (m *MockExecutor) CommandContext(ctx context.Context, name string, args ...string) error {
    args = append([]string{name}, args...)
    argsInterface := make([]interface{}, len(args))
    for i, v := range args {
        argsInterface[i] = v
    }
    ret := m.Called(append([]interface{}{ctx}, argsInterface...)...)
    return ret.Error(0)
}

func (m *MockExecutor) Command(name string, args ...string) (string, error) {
    args = append([]string{name}, args...)
    argsInterface := make([]interface{}, len(args))
    for i, v := range args {
        argsInterface[i] = v
    }
    ret := m.Called(argsInterface...)
    return ret.String(0), ret.Error(1)
}

func (m *MockExecutor) GH(args ...string) (bytes.Buffer, error) {
    argsInterface := make([]interface{}, len(args))
    for i, v := range args {
        argsInterface[i] = v
    }
    ret := m.Called(argsInterface...)
    return ret.Get(0).(bytes.Buffer), ret.Error(1)
}

func TestStatusDummy(t *testing.T) {
    mockExec := &MockExecutor{}
    statusChecker := status.NewStatus(mockExec)

    result, err := statusChecker.GetStatus("dummyArgument")

    assert.NoError(t, err)
    assert.NotNil(t, result)
}