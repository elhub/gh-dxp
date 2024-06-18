package lint_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/lint"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockExecutor struct {
	mock.Mock
}

func (m *mockExecutor) Command(name string, arg ...string) (string, error) {
	args := m.Called(name, arg)
	return args.String(0), args.Error(1)
}

func (m *mockExecutor) CommandContext(ctx context.Context, name string, arg ...string) error {
	args := m.Called(ctx, name, arg)
	return args.Error(1)
}

func (m *mockExecutor) GH(arg ...string) (bytes.Buffer, error) {
	args := m.Called(arg)
	return *bytes.NewBufferString(args.String(0)), args.Error(1)
}

func TestRun(t *testing.T) {
	tests := []struct {
		name      string
		mockError error
		expectErr bool
	}{
		{
			name:      "command runs successfully",
			mockError: nil,
			expectErr: false,
		},
		{
			name:      "command returns an error",
			mockError: errors.New("command error"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExe := new(mockExecutor)
			args := []string{"mega-linter-runner", "--flavor", "cupcake", "-e",
				"MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml"}
			mockExe.On("CommandContext", mock.Anything, "npx", args).Return(nil, tt.mockError)

			err := lint.Run(mockExe, &config.Settings{})

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			mockExe.AssertExpectations(t)
		})
	}
}
