package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/merge"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

func MergeCmd(exe utils.Executor, settings *config.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge",
		Short: "Merge a pull request",
		Long: heredoc.Docf(`
			Merge the pull request on the current branch. This is an opinionated command that will:

			* Squash-merge the current branch to main or master
			* Delete both the local and remote branches
		`, "`"),
		Example: heredoc.Doc(`
			$ gh dxp merge
		`),
		Aliases: []string{"land"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return merge.Execute(exe)
		},
	}

	// TODO: Support flags from gh pr

	return cmd
}
