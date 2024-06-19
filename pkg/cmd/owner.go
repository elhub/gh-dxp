package cmd

import (
	"errors"

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
		Args:  cobra.MaximumNArgs(1),
		Long: heredoc.Docf(`
			Determine the owner of the specified file based on the CODEOWNERS file given in the .github directory. If
			no CODEOWNERS is found in the .github directory, it will return undefined.
		`, "`"),
		Example: heredoc.Doc(`
			// Check the owner of the README.md file
			$ gh dxp owner README.md
		`),
		RunE: func(_ *cobra.Command, args []string) error {
			var path string
			if len(args) > 1 {
				path = args[0]
			} else {
				rootDir, err := utils.GetGitRootDirectory(exe)
				if err != nil {
					return err
				}

				readmeFile := rootDir + "README.md"

				if utils.FileExists(readmeFile) {
					path = readmeFile
				} else {
					// Get the first file in the root directory
					files, err := utils.ListFilesInDirectory(exe, rootDir)
					if err != nil {
						return err
					}
					if len(files) > 0 {
						path = rootDir + files[0]
					} else {
						return errors.New("no files found in the root directory")
					}

				}
			}

			owners, err := owner.Execute(path, exe)

			// Output the owners
			for _, owner := range owners {
				println(owner)
			}

			return err
		},
	}

	return cmd
}
