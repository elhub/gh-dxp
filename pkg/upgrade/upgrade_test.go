package upgrade_test

import (
	"errors"
	"testing"

	"github.com/elhub/gh-dxp/pkg/testutils"
	"github.com/elhub/gh-dxp/pkg/upgrade"
	"github.com/stretchr/testify/require"
)

func TestRunUpgrade(t *testing.T) {
	ghArgs := []string{"extension", "upgrade", "elhub/gh-dxp", "--force"}

	t.Run("Upgrade successful", func(t *testing.T) {
		mockExe := new(testutils.MockExecutor)
		mockExe.On("GH", ghArgs).Return("", nil)

		err := upgrade.RunUpgrade(mockExe)
		require.NoError(t, err)
	})
	t.Run("Upgrade unsuccessful", func(t *testing.T) {
		mockExe := new(testutils.MockExecutor)
		mockExe.On("GH", ghArgs).Return("", errors.New("something went wrong"))

		err := upgrade.RunUpgrade(mockExe)
		require.Error(t, err)
	})
}
