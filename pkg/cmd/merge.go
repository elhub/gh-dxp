package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/elhub/gh-dxp/pkg/merge"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

// MergeCmd handles the merge of a pull request.
func MergeCmd(exe utils.Executor) *cobra.Command {
	opts := &merge.Options{}

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
		RunE: func(_ *cobra.Command, _ []string) error {
			return merge.Execute(exe, opts)
		},
	}

	// TODO: Support flags from gh pr
	fl := cmd.Flags()
	fl.BoolVarP(
		&opts.AutoConfirm,
		"confirm",
		"y",
		false,
		"Don't ask for user input.",
	)

	return cmd
}
