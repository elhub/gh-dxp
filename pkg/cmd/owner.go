package cmd

import (
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/owner"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

// OwnerCmd creates a new cobra command for retrieving code owner information.
func OwnerCmd(exe utils.Executor, _ *config.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner",
		Short: "Determines the owner of the specified file.",
		Args:  cobra.MaximumNArgs(1),
		Long: heredoc.Docf(`
			Determine the owner of the specified file based on the CODEOWNERS file given in the .github directory. If
			no CODEOWNERS is found in the .github directory, it will return undefined.
		`, "`"),
		Example: heredoc.Doc(`
			# Check the owner of the README.md file
			$ gh dxp owner README.md
		`),
		RunE: func(_ *cobra.Command, args []string) error {
			var path string
			if len(args) > 0 {
				path = args[0]
			} else {
				defaultPath, err := owner.GetDefaultFile(exe)
				if err != nil {
					return err
				}

				path = defaultPath
			}

			owners, err := owner.Execute(path, exe)

			// Output the owners
			for _, owner := range owners {
				_, err := os.Stdout.Write([]byte(owner + "\n"))
				if err != nil {
					return err
				}
			}

			return err
		},
	}

	return cmd
}
