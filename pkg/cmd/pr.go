package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/pr"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

func PRCmd(exe utils.Executor, settings *config.Settings) *cobra.Command {
	opts := &pr.Options{}

	cmd := &cobra.Command{
		Use:   "pr",
		Short: "Create a PR (pull request)",
		Long: heredoc.Docf(`
			Create a PR (pull request) for the current branch. This is an opinionated command that will:

			* Push the current branch to git remote
			* Generate a pull request title based on the current branch name (based upon settings)
			* Generate a pull request body based on the current devxp template
		`, "`"),
		Example: heredoc.Doc(`
			$ gh dxp pr
		`),
		Aliases: []string{"diff"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return pr.Execute(exe, opts)
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
	fl.StringSliceVarP(
		&opts.Reviewers,
		"reviewer",
		"r",
		nil,
		"Request reviews from people or teams by their id",
	)
	fl.StringSliceVarP(
		&opts.Assignees,
		"assignee",
		"a",
		nil,
		"Assign people by their id. Use \"@me\" to self-assign.",
	)

	return cmd
}
