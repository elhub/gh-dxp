package testutils_test

import (
	"testing"

	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMockExecutor(t *testing.T) {
	mocks := []testutils.MockContent{
		{
			Method: "Command",
			Args:   []interface{}{"git", []string{"rev-parse", "--show-toplevel"}},
			Out:    "/path/to/repo\n",
			Err:    nil,
		},
		{
			Method: "CommandContext",
			Args:   []interface{}{nil, "git", []string{"rev-parse", "--show-toplevel"}},
			Out:    "",
			Err:    nil,
		},
		{
			Method: "GH",
			Args:   []interface{}{[]string{"repo", "view"}},
			Out:    "mocked output",
			Err:    nil,
		},
		{
			Method: "Chdir",
			Args:   []interface{}{"/mock/path"},
			Out:    "",
			Err:    nil,
		},
	}

	mockExe := testutils.NewMockExecutor(mocks)

	// Test Command method
	cOutput, err := mockExe.Command("git", "rev-parse", "--show-toplevel")
	assert.Equal(t, "/path/to/repo\n", cOutput)
	require.NoError(t, err)

	// Test CommandContext method
	err = mockExe.CommandContext(nil, "git", "rev-parse", "--show-toplevel")
	require.NoError(t, err)

	// Test GH method
	ghOutput, err := mockExe.GH("repo", "view")
	require.NoError(t, err)
	assert.Equal(t, "mocked output", ghOutput)

	// Test Chdir method
	err = mockExe.Chdir("/mock/path")
	require.NoError(t, err)

	mockExe.AssertExpectations(t)
}
