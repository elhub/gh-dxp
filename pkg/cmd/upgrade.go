package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/elhub/gh-dxp/pkg/upgrade"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

// UpgradeCmd upgrades the dxp plugin to the latest version
func UpgradeCmd(exe utils.Executor) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "upgrade",
		Short:  "Upgrade gh dxp",
		Hidden: true,
		Args:   cobra.ExactArgs(0),
		Long: heredoc.Docf(`
			Upgrade gh dxp to latest version`, "`"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return upgrade.RunUpgrade(exe)
		},
	}
	return cmd
}
