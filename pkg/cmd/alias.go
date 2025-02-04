package cmd

import (
	"context"
	"os/exec"

	"github.com/MakeNowJust/heredoc"
	"github.com/elhub/gh-dxp/pkg/logger"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// AliasCmd extends the functionality of the gh alias command to facilitate the import of Elhub's preferred
// aliases.
func AliasCmd(exe utils.Executor) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alias",
		Short: "Import default aliases.",
		Args:  cobra.ExactArgs(1),
		Long: heredoc.Docf(`
			Aliases are used to create shortcuts for gh commands. This command allows to import the default
			aliases used by Elhub.
		`),
	}

	cmd.AddCommand(AliasImportCmd(exe))

	return cmd
}

// AliasImportCmd imports the default aliases from github.
func AliasImportCmd(exe utils.Executor) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import default aliases.",
		Long: heredoc.Docf(`
			Import the default aliases used by Elhub (stored in elhub/gh-dxp as alias.yml).
		`, "`"),
		Example: heredoc.Doc(`
			# Import the default aliases
			$ gh dxp alias import
		`),
		RunE: func(_ *cobra.Command, _ []string) error {
			logger.Info("Piping default aliases from github/elhub/gh-dxp to standard input")
			ctx := context.Background()
			err := exe.CommandContext(ctx, "sh", "-c",
				"curl -s https://raw.githubusercontent.com/elhub/gh-dxp/main/alias.yml | gh alias import - --clobber")

			if err != nil {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					if exitErr.ExitCode() == 0 {
						return nil
					}
					return nil
				}

				return err
			}
			return nil
		},
	}

	return cmd
}
