package cmd

import (
	"github.com/MakeNowJust/heredoc"
	merge "github.com/elhub/gh-dxp/pkg/prmerge"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

// MergeCmd handles the merge of a pull request.
func PRMergeCmd(exe utils.Executor) *cobra.Command {
	opts := &merge.Options{}

	cmd := &cobra.Command{
		Use:   "pr-merge",
		Short: "Merge a PR (Pull Request)",
		Long: heredoc.Docf(`
			Merge the pull request on the current branch. This is an opinionated command that will:

			* Squash-merge the current branch to main or master
			* Delete both the local and remote branches
		`, "`"),
		Example: heredoc.Doc(`
			# Merge the current branch if it is a GitHub PR
			$ gh dxp pr-merge
		`),
		Aliases: []string{"land", "merge"},
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			err := utils.SetWorkDirToGitHubRoot(exe)
			if err != nil {
				return err
			}

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