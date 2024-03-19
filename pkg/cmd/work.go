package cmd

import (
	"context"
	"os/exec"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func WorkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "work [branch-name]",
		Short: "Create a new branch based on an issue and checkout to it.",
		Args:  cobra.MaximumNArgs(1),
		Long: heredoc.Docf(`
			Create a new branch and checkout to it. If the branch already exists,
			it will be checked out.
		`, "`"),
		Example: heredoc.Doc(`
			// Create a new branch 'wip' and checkout to it:
			$ gh devxp work wip
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			branchID := ""
			if len(args) > 0 {
				branchID = args[0]
			}

			return checkoutBranch(ctx, branchID)
		},
	}

	return cmd
}

func checkoutBranch(_ context.Context, branchID string) error {
	// Does the branch exist?
	branchExists, existsErr := branchExists(branchID)
	if existsErr != nil {
		return existsErr
	}

	if branchExists {
		log.Debugf("Branch '%s' already exists, checking out to it", branchID)
		out, err1 := utils.Exec().Run("git", "checkout", branchID)
		if err1 != nil {
			return errors.Wrap(err1, "Failed to checkout branch")
		}
		log.Info(strings.Trim(out, "\n"))
	} else {
		log.Debugf("Creating branch '%s' and checking out to it", branchID)
		out, err2 := utils.Exec().Run("git", "checkout", "-b", branchID)
		if err2 != nil {
			return errors.Wrap(err2, "Failed to create branch")
		}
		log.Info(strings.Trim(out, "\n"))
	}

	return nil
}

func branchExists(branchID string) (bool, error) {
	_, err := utils.Exec().Run("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branchID)
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if exitErr.ExitCode() == 0 {
				return true, nil
			}
			return false, nil
		}

		return false, err
	}
	return true, nil
}
