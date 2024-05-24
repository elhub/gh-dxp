package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/elhub/gh-dxp/pkg/branch"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

func BranchCmd(exe utils.Executor) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "branch [branch-name]",
		Short: "Create a new branch based on an issue and checkout to it.",
		Args:  cobra.ExactArgs(1),
		Long: heredoc.Docf(`
			Create a new branch and checkout to it. If the branch already exists,
			it will be checked out.
		`, "`"),
		Example: heredoc.Doc(`
			// Create a new branch 'wip' and checkout to it:
			$ gh devxp branch wip
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			branchID := args[0]

			s := utils.StartSpinner("Creating new work branch...", "Work Branch "+branchID)
			b := branch.CheckoutBranch(exe, branchID)
			s.Stop()
			return b
		},
	}

	return cmd
}
