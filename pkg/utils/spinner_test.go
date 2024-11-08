package utils_test

import (
	"fmt"
	"testing"

	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestStartSpinner(t *testing.T) {
	suffix := "Suffix"
	finalMsg := "Final Message"

	t.Run("Spinner is created", func(t *testing.T) {

		s := utils.StartSpinner(suffix, finalMsg)

		assert.Equal(t, fmt.Sprintf(" %s", suffix), s.Suffix)
		assert.Equal(t, fmt.Sprintf(" %s\n", finalMsg), s.FinalMSG)

	})

	t.Run("Spinner final message is correctly removed", func(t *testing.T) {

		s := utils.StartSpinner(suffix, finalMsg)
		utils.RemoveSpinnerFinalMsg(s)

		assert.Empty(t, s.FinalMSG)
	})
}
