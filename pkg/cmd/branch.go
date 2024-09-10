package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/elhub/gh-dxp/pkg/branch"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

// BranchCmd creates a new branch based on an issue and checks out to it.
func BranchCmd(exe utils.Executor) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "branch [branch-name]",
		Short: "Checkout or Create+Checkout a git branch.",
		Args:  cobra.ExactArgs(1),
		Long: heredoc.Docf(`
			Create a new branch and checkout to it. If the branch already exists,
			it will be checked out.
		`, "`"),
		Example: heredoc.Doc(`
			# Create a new branch 'wip' and checkout to it:
			$ gh dxp branch wip
		`),
		RunE: func(_ *cobra.Command, args []string) error {
			err := utils.SetWorkDirToGitHubRoot(exe)
			if err != nil {
				return err
			}
			branchID := args[0]

			s := utils.StartSpinner("Creating new work branch...", "Work Branch "+branchID)
			b := branch.CheckoutBranch(exe, branchID)
			s.Stop()
			return b
		},
	}

	return cmd
}
