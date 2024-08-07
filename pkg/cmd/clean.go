package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/elhub/gh-dxp/pkg/clean"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

// BranchCmd creates a new branch based on an issue and checks out to it.
func CleanCmd(exe utils.Executor) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean [target]",
		Short: "Clean repo of excess files or branches",
		Args:  cobra.ExactArgs(1),
		Long: heredoc.Docf(`
				Branches: Delete all local branches except for the default branch.
				Files: Delete any untracked or gitignored files in the repository.
		`, "`"),
		Example: heredoc.Doc(`
			// Delete all non-default branches
			$ gh dxp clean branch
			// Delete all untracked and gitignored files 
			$ gh dxp clean file
		`),
		RunE: func(_ *cobra.Command, args []string) error {
			target := args[0]

			s := utils.StartSpinner("Cleaning repository...", "Target "+target)
			b := clean.Run(exe, target)
			s.Stop()
			return b
		},
	}

	return cmd
}
