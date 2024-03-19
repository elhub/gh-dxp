package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/elhub/gh-devxp/pkg/config"
	"github.com/elhub/gh-devxp/pkg/lint"
	"github.com/spf13/cobra"
)

func LintCmd(settings *config.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Run the set of linters defined in the .devxp config.",
		Args:  cobra.MaximumNArgs(0),
		Long: heredoc.Docf(`
			Run the set of linters defined in the .devxp config file. If no linters are defined,
			the command will do nothing.
		`, "`"),
		Example: heredoc.Doc(`
			// Lint the current directory
			$ gh devxp lint
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			return lint.Run(ctx, settings, lint.DefaultLinters())
		},
	}

	return cmd
}
