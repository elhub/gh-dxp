package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/pr"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

// PRCmd extends the functionality of the gh pr command.
func PRCmd(exe utils.Executor, settings *config.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pr",
		Short: "Work with PRs (Pull Requests)",
		Long: heredoc.Doc(`
			The pr command group allows you to create, view, and manage pull requests. It extends the
			functionality of the gh pr command.
		`),
	}

	cmd.AddCommand(PRCreateCmd(exe, settings))
	cmd.AddCommand(PRMergeCmd(exe))
	cmd.AddCommand(PRUpdateCmd(exe, settings))

	var opts = &pr.Options{}

	fl := cmd.PersistentFlags()
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
	fl.StringVarP(
		&opts.CommitMessage,
		"commitmessage",
		"m",
		"",
		"Commit message, if there are uncommitted changes.",
	)
	return cmd
}

// PRCreateCmd handles the creation of a pull request.
func PRCreateCmd(exe utils.Executor, settings *config.Settings) *cobra.Command {

	opts := &pr.CreateOptions{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a PR (Pull Request)",
		Long: heredoc.Docf(`
			Create a PR (pull request) from the current branch. This is an opinionated command
			that will:

			* Run tests and linting
			* Push the current branch to the remote repository
			* Create a pull request on GitHub
		`, "`"),
		Example: heredoc.Doc(`
			# Create a PR from the current branch
			$ gh dxp pr create
		`),
		Args: cobra.NoArgs,
		RunE: func(prCmd *cobra.Command, _ []string) error {
			prOpts, err := getPrOptionsFromCmd(prCmd)
			if err != nil {
				return err
			}

			addPrOptionsToCreateOptions(prOpts, opts)

			err = utils.SetWorkDirToGitHubRoot(exe)
			if err != nil {
				return err
			}
			return pr.ExecuteCreate(exe, settings, opts)
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
		&opts.Branch,
		"branch",
		"b",
		"",
		"Temporary branch to switch to if currently on the default branch",
	)
	fl.StringVarP(
		&opts.Issues,
		"issues",
		"i",
		"",
		"Comma separated list of issues",
	)
	fl.BoolVar(
		&opts.Draft,
		"draft",
		false,
		"Mark pull request as draft",
	)

	return cmd
}

// PRMergeCmd handles the merging of a pull request.
func PRMergeCmd(exe utils.Executor) *cobra.Command {
	opts := &pr.MergeOptions{}

	cmd := &cobra.Command{
		Use:   "merge",
		Short: "Merge a PR (Pull Request)",
		Long: heredoc.Docf(`
			Merge the pull request on the current branch. This is an opinionated command that will:

			* Squash-merge the current branch to main or master
			* Delete both the local and remote branches
		`, "`"),
		Example: heredoc.Doc(`
			# Merge the current branch if it is a GitHub PR
			$ gh dxp pr merge
		`),
		Aliases: []string{"land", "merge"},
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			err := utils.SetWorkDirToGitHubRoot(exe)
			if err != nil {
				return err
			}

			return pr.ExecuteMerge(exe, opts)
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

// PRUpdateCmd handles the updating of a pull request. This is a more limited version of the create command.
func PRUpdateCmd(exe utils.Executor, settings *config.Settings) *cobra.Command {
	opts := &pr.UpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a PR (Pull Request)",
		Long: heredoc.Docf(`
			Update the current PR (pull request). This is essentially the same as the
			create command, except that it expects an existing PR.
		`, "`"),
		Example: heredoc.Doc(`
			# Update the current PR
			$ gh dxp pr update
		`),
		Args: cobra.NoArgs,
		RunE: func(prCmd *cobra.Command, _ []string) error {
			prOpts, err := getPrOptionsFromCmd(prCmd)
			if err != nil {
				return err
			}

			addPrOptionsToUpdateOptions(prOpts, opts)
			err = utils.SetWorkDirToGitHubRoot(exe)
			if err != nil {
				return err
			}
			return pr.ExecuteUpdate(exe, settings, opts)
		},
	}

	return cmd
}

// Patches prOptions into a createOptions object
func addPrOptionsToCreateOptions(prOptions pr.Options, createOptions *pr.CreateOptions) {
	createOptions.NoLint = prOptions.NoLint
	createOptions.NoUnit = prOptions.NoUnit
	createOptions.CommitMessage = prOptions.CommitMessage
}

// Patches prOptions into an updateOptions object
func addPrOptionsToUpdateOptions(prOptions pr.Options, updateOptions *pr.UpdateOptions) {
	updateOptions.NoLint = prOptions.NoLint
	updateOptions.NoUnit = prOptions.NoUnit
	updateOptions.CommitMessage = prOptions.CommitMessage
}

func getPrOptionsFromCmd(cmd *cobra.Command) (pr.Options, error) {
	var prOptions pr.Options

	noLint, err := cmd.Flags().GetBool("nolint")
	if err != nil {
		return pr.Options{}, err
	}
	noUnit, err := cmd.Flags().GetBool("nounit")
	if err != nil {
		return pr.Options{}, err
	}
	commitMessage, err := cmd.Flags().GetString("commitmessage")
	if err != nil {
		return pr.Options{}, err
	}
	prOptions.NoLint = noLint
	prOptions.NoUnit = noUnit
	prOptions.CommitMessage = commitMessage
	return prOptions, nil
}
