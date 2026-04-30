package cmd

import (
	"github.com/elhub/gh-dxp/pkg/pr"
	"github.com/spf13/cobra"
)

var AddPrOptionsToCreateOptions = addPrOptionsToCreateOptions //nolint:gochecknoglobals // Expose for testing
var AddPrOptionsToUpdateOptions = addPrOptionsToUpdateOptions //nolint:gochecknoglobals // Expose for testing

func GetPrOptionsFromCmd(cmd *cobra.Command) (pr.Options, error) {
	return getPrOptionsFromCmd(cmd)
}
