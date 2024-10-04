package utils_test

import (
	"errors"
	"testing"

	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLatestReleaseVersion(t *testing.T) {

	ghAPIArgs := []string{"api", "-H", "Accept: application/vnd.github+json", "-H", "X-GitHub-Api-Version: 2022-11-28", "/repos/elhub/gh-dxp/releases/latest"}

	t.Run("Latest release version is v1.2.3", func(t *testing.T) {
		mockExe := new(testutils.MockExecutor)
		mockExe.On("GH", ghAPIArgs).Return(`{"tag_name":"v1.2.3"}`, nil)

		version, err := utils.GetLatestReleaseVersion(mockExe)

		require.NoError(t, err)
		assert.Equal(t, "v1.2.3", version)
	})

	t.Run("GitHub response is incorrectly formatted", func(t *testing.T) {
		mockExe := new(testutils.MockExecutor)
		mockExe.On("GH", ghAPIArgs).Return(`{"tag_name"v1.2.3"}`, nil)

		_, err := utils.GetLatestReleaseVersion(mockExe)

		require.Error(t, err)
	})
	t.Run("No tag_name in response", func(t *testing.T) {
		mockExe := new(testutils.MockExecutor)
		mockExe.On("GH", ghAPIArgs).Return(`{"other_field":"v1.2.3"}`, nil)

		_, err := utils.GetLatestReleaseVersion(mockExe)

		require.Error(t, err)
	})

}

func TestIsLatestVersionOrSnapshot(t *testing.T) {
	ghAPIArgs := []string{"api", "-H", "Accept: application/vnd.github+json", "-H", "X-GitHub-Api-Version: 2022-11-28", "/repos/elhub/gh-dxp/releases/latest"}

	t.Run("Local version is SNAPSHOT", func(t *testing.T) {
		mockExe := new(testutils.MockExecutor)

		result, err := utils.IsLatestVersionOrSnapshot(mockExe, "SNAPSHOT")

		require.NoError(t, err)
		assert.Equal(t, true, result)
	})
	t.Run("Local version is not latest", func(t *testing.T) {
		mockExe := new(testutils.MockExecutor)
		mockExe.On("GH", ghAPIArgs).Return(`{"tag_name":"v1.2.3"}`, nil)

		result, err := utils.IsLatestVersionOrSnapshot(mockExe, "v1.2.2")

		assert.Equal(t, false, result)
		require.NoError(t, err)

	})
	t.Run("Local version is latest", func(t *testing.T) {
		mockExe := new(testutils.MockExecutor)
		mockExe.On("GH", ghAPIArgs).Return(`{"tag_name":"v1.2.3"}`, nil)

		result, err := utils.IsLatestVersionOrSnapshot(mockExe, "v1.2.3")

		assert.Equal(t, true, result)
		require.NoError(t, err)

	})
	t.Run("Error in api call", func(t *testing.T) {
		mockExe := new(testutils.MockExecutor)
		mockExe.On("GH", ghAPIArgs).Return("", errors.New("Some error"))

		_, err := utils.IsLatestVersionOrSnapshot(mockExe, "v1.2.3")

		require.Error(t, err)

	})

}
