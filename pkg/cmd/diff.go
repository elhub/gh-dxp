package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/caarlos0/log"
	"github.com/cli/go-gh"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type DiffOptions struct {
	Confirm bool

	baseBranch string

	Reviewers []string
	Assignees []string
}

func DiffCmd(settings *config.Settings) *cobra.Command {
	opts := &DiffOptions{}

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Create a diff (pull request)",
		Long: heredoc.Docf(`
			Create a diff (pull request) for the current branch. This is an opinionated command that will:

			* Push the current branch to git remote
			* Generate a pull request title based on the current branch name (based upon settings)
			* Generate a pull request body based on the current devxp template
		`, "`"),
		Example: heredoc.Doc(`
			$ gh dxp diff
		`),
		Aliases: []string{"create", "pr", "new"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createDiff(cmd.Context(), settings, opts)
		},
	}
	// TODO: Support flags grom gh pr

	fl := cmd.Flags()
	fl.BoolVarP(
		&opts.Confirm,
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

func createDiff(_ context.Context, settings config.Settings, options *DiffOptions) error {
	// Set branchID
	currentBranch, err := utils.Exec().Run("git", "branch", "--show-current")
	if err != nil {
		return err
	}
	branchId := strings.Trim(currentBranch, "\n")

	// Push the current branch to git remote
	//s := utils.StartSpinner("Pushing current branch to remote...", "Pushed branch to remote")
	currentBranch, err = utils.Exec().Run("git", "push", "--set-upstream", "origin", branchId)
	//s.Stop()
	if err != nil {
		return err
	}
	log.Info(strings.Trim(currentBranch, "\n"))

	// Fetch the current branch
	baseBranch := options.baseBranch
	if baseBranch == "" {
		//s := utils.StartSpinner("Fetching repository default branch...", "Fetched repository default branch")
		stdOut, _, err := gh.Exec("repo", "view", "--json", "defaultBranchRef", "--jq", ".defaultBranchRef.name")
		//s.Stop()
		if err != nil {
			return errors.Wrap(err, "Failed to fetch default branch")
		}
		baseBranch = strings.Trim(stdOut.String(), "\n")
	}

	pr, err := pr.createPR(currentBranch, options.Confirm)
	if err != nil {
		return err
	}

	log.Debug(fmt.Sprintf("Pull request title: %s", pr.Title))
	log.Debug(fmt.Sprintf("Pull request body:\n\n%s", pr.Body))

	//s := utils.StartSpinner("Creating pull request...", "Created pull request")
	args := []string{"pr", "create", "--title", pr.Title, "--body", pr.Body, "--base", baseBranch}
	args = append(args, generatePRArgs(options)...)
	stdOut, _, err := gh.Exec(args...)
	//s.Stop()
	if err != nil {
		return errors.Wrap(err, "Failed to create pull request")
	}
	log.Info(strings.Trim(stdOut.String(), "\n"))

	return nil
}

func generatePRArgs(options *DiffOptions) []string {
	args := []string{}

	if len(options.Assignees) > 0 {
		args = append(args, "--assignee", strings.Join(options.Assignees, ","))
	}
	if len(options.Reviewers) > 0 {
		args = append(args, "--reviewer", strings.Join(options.Reviewers, ","))
	}

	return args
}
