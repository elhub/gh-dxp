package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/lint"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

// LintCmd creates a new command to run the linters defined in the .devxp config.
func LintCmd(exe utils.Executor, settings *config.Settings) *cobra.Command {
	opts := &lint.Options{}

	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Run MegaLinter on modified files in the repository.",
		Args:  cobra.MaximumNArgs(0),
		Long: heredoc.Docf(`
			Run linters on files in the repository. By default, only files that have been modified in relation to the main branch are included in the lint.
		`, "`"),
		Example: heredoc.Doc(`
			// Lint modified files in the repository
			$ gh dxp lint

			// Lint all files in repository
			$ gh dxp lint --all
		`),
		RunE: func(_ *cobra.Command, _ []string) error {
			return lint.Run(exe, settings, opts)
		},
	}

	fl := cmd.Flags()
	fl.BoolVarP(
		&opts.LintAll,
		"all",
		"a",
		false,
		"Lint all files in the repository",
	)

	return cmd
}
