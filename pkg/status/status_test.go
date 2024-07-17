package status

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExecutor is a mock type for the Executor type
type MockExecutor struct {
	mock.Mock
}

// CommandContext implements utils.Executor.
func (m *MockExecutor) CommandContext(ctx context.Context, name string, arg ...string) error {
	panic("unimplemented")
}

// GH implements utils.Executor.
func (m *MockExecutor) GH(args ...string) (bytes.Buffer, error) {
	panic("unimplemented")
}

// Command is a mock method that simulates the behavior of utils.Executor's Command method
func (m *MockExecutor) Command(name string, arg ...string) (string, error) {
	args := m.Called(name, arg)
	return args.String(0), args.Error(1)
}

func TestGetStatus(t *testing.T) {
	tests := []struct {
		name          string
		executorSetup func(*MockExecutor)
		want          string
		wantErr       bool
	}{
		{
			name: "successful status retrieval",
			executorSetup: func(m *MockExecutor) {
				m.On("Command", "git", []string{"status"}).Return("On branch main\nYour branch is up to date with 'origin/main'.", nil)
			},
			want:    "On branch main\nYour branch is up to date with 'origin/main'.",
			wantErr: false,
		},
		{
			name: "error retrieving status",
			executorSetup: func(m *MockExecutor) {
				m.On("Command", "git", []string{"status"}).Return("", errors.New("error executing git status"))
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExecutor := new(MockExecutor)
			tt.executorSetup(mockExecutor)

			s := NewStatus(mockExecutor)
			got, err := s.GetStatus()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
