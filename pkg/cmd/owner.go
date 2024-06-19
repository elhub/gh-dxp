package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/owner"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

// LintCmd creates a new command to run the linters defined in the .devxp config.
func OwnerCmd(exe utils.Executor, settings *config.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner",
		Short: "Determines the owner of the specified file.",
		Args:  cobra.ExactArgs(1),
		Long: heredoc.Docf(`
			Determine the owner of the specified file based on the CODEOWNERS file given in the .github directory. If
			no CODEOWNERS is found in the .github directory, it will return undefined.
		`, "`"),
		Example: heredoc.Doc(`
			// Check the owner of the README.md file
			$ gh dxp owner README.md
		`),
		RunE: func(_ *cobra.Command, args []string) error {
			path := args[0]
			return owner.Execute(path)
		},
	}

	return cmd
}
