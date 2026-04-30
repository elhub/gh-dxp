// Package cmd provides CLI commands for the gh-dxp extension.
package cmd

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/elhub/gh-dxp/pkg/ghutil"
	"github.com/elhub/gh-dxp/pkg/test"
	"github.com/spf13/cobra"
)

// TestCmd handles the running of tests.
func TestCmd(exe ghutil.Executor) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Run tests",
		Args:  cobra.ExactArgs(0),
		Long: heredoc.Docf(`
			Run tests based on project type`, "`"),
		RunE: func(_ *cobra.Command, _ []string) error {
			err := ghutil.SetWorkDirToGitHubRoot(exe)
			if err != nil {
				return err
			}
			_, err = test.RunTest(exe)
			return err
		},
	}
	return cmd
}
