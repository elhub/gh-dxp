package cmd

import (
	"github.com/elhub/gh-dxp/pkg/pr"
	"github.com/spf13/cobra"
)

var AddPrOptionsToCreateOptions = addPrOptionsToCreateOptions
var AddPrOptionsToUpdateOptions = addPrOptionsToUpdateOptions

func GetPrOptionsFromCmd(cmd *cobra.Command) (pr.Options, error) {
	return getPrOptionsFromCmd(cmd)
}
