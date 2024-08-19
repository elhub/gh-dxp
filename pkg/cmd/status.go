package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/elhub/gh-dxp/pkg/status"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

// StatusCmd creates a new cobra.Command for the status functionality.
func StatusCmd(exe utils.Executor) *cobra.Command {
	opts := &status.Options{}

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Get the status of the current repository",
		Long: heredoc.Docf(`
            Get the status of the current repository or specific aspects related to it. This command allows you to:

            * View a comprehensive status report covering all aspects (All)
            * Check the status of the current repository
            * View the status of pull requests (PRs)
            * List all branches within the repository
            * View assigned PRs/Review Requests

            This command supports both interactive mode and non-interactive mode via flags for quick access to specific information.
        `, "`"),
		Example: heredoc.Doc(`
            # Interactive mode
            $ gh dxp status
        `),
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			err := utils.SetWorkDirToGitHubRoot(exe)
			if err != nil {
				return err
			}
			return status.Execute(exe, opts)
		},
	}

	fl := cmd.Flags()
	fl.BoolVarP(
		&opts.All,
		"all",
		"a",
		false,
		"Get all statuses",
	)
	fl.BoolVarP(
		&opts.Repo,
		"current",
		"c",
		false,
		"Get current repository target",
	)
	fl.BoolVarP(
		&opts.Pr,
		"pr",
		"p",
		false,
		"Get the PR status",
	)
	fl.BoolVarP(
		&opts.Branches,
		"branches",
		"b",
		false,
		"List all branches",
	)
	fl.BoolVarP(
		&opts.Issue,
		"issues",
		"i",
		false,
		"Get all relavant Issues",
	)

	return cmd
}
