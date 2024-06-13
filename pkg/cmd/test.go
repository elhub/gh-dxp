package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/elhub/gh-dxp/pkg/test"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

// TestCmd handles the running of tests.
func TestCmd(exe utils.Executor) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Run tests",
		Args:  cobra.ExactArgs(0),
		Long: heredoc.Docf(`
			Run tests based on project type`, "`"),
		RunE: func(_ *cobra.Command, _ []string) error {
			res := test.RunTest(exe)
			return res
		},
	}
	return cmd
}
