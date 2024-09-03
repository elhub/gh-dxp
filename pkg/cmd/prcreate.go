package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/prcreate"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

// PRCreateCmd handles the creation of a pull request.
func PRCreateCmd(exe utils.Executor, settings *config.Settings) *cobra.Command {
	opts := &prcreate.Options{}

	cmd := &cobra.Command{
		Use:   "pr",
		Short: "Create a PR (Pull Request)",
		Long: heredoc.Docf(`
			Create a PR (pull request) for the current branch. This is an opinionated command that will:

			* Push the current branch to git remote
			* Generate a pull request title based on the current branch name (based upon settings)
			* Generate a pull request body based on the current devxp template
		`, "`"),
		Example: heredoc.Doc(`
			# Create a PR from the current branch
			$ gh dxp pr-create
		`),
		Aliases: []string{"diff", "pr"},
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			err := utils.SetWorkDirToGitHubRoot(exe)
			if err != nil {
				return err
			}
			return prcreate.Execute(exe, settings, opts)
		},
	}

	// TODO: Support flags from gh pr
	fl := cmd.Flags()
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
	fl.StringVarP(
		&opts.CommitMessage,
		"commitmessage",
		"m",
		"",
		"Commit message, if there are uncommitted changes.",
	)
	fl.StringVarP(
		&opts.Branch,
		"branch",
		"b",
		"",
		"Temporary branch to switch to if currently on the default branch",
	)
	fl.BoolVar(
		&opts.NoUnit,
		"nounit",
		false,
		"Do not run tests",
	)
	fl.BoolVar(
		&opts.NoLint,
		"nolint",
		false,
		"Do not run linting",
	)

	return cmd
}
