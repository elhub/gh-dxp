package cmd

import (
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/repo"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

// RepoCmd extends the functionality of the gh repo command.
func RepoCmd(exe utils.Executor, _ *config.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo",
		Short: "Work with Repositories",
		Long: heredoc.Doc(`
			The repo command group extends the basic repo commands provided by the
			standard GitHub CLI commands with some extra bells and whistles useful for
			handling a set of repositories.
		`),
	}

	cmd.AddCommand(RepoCloneCmd(exe))

	return cmd
}

// RepoCloneCmd creates a new command to clone all repositories (or just those with a given name).
func RepoCloneCmd(exe utils.Executor) *cobra.Command {
	opts := &repo.Options{}

	cmd := &cobra.Command{
		Use:   "clone-all [<pattern>]",
		Short: "Clone all repositories with a given name.",
		Long: heredoc.Docf(`
			Clones all repositories with a given name pattern. If no pattern is provided,
			clone all repositories from the user's organizations.
		`, "`"),
		Example: heredoc.Doc(`
			# Clone everything
			$ gh dxp repo clone-all

			# Clone all repositories in any of your organizations starting with gh
			$ gh dxp repo clone-all gh

			# Check what repositories would be cloned
			$ gh dxp repo clone-all gh --dryrun
		`),
		RunE: func(_ *cobra.Command, args []string) error {
			pattern := ""
			if len(args) > 0 {
				pattern = args[0]
			}
			sleepFunction := time.Sleep
			return repo.ExecuteClone(exe, pattern, sleepFunction, opts)
		},
	}

	fl := cmd.Flags()
	fl.BoolVar(
		&opts.DryRun,
		"dryrun",
		false,
		"Do not actually clone the repositories, just list the repositories that would be cloned.",
	)

	return cmd
}
