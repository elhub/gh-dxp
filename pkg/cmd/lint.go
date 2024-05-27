package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/lint"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

func LintCmd(exe utils.Executor, settings *config.Settings) *cobra.Command {
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
			$ gh dxp lint
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return lint.Run(exe, settings)
		},
	}

	return cmd
}
