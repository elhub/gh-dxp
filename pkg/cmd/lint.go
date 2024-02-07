package cmd

import (
	"context"
	"os/exec"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/caarlos0/log"
	"github.com/michaeloa/gh-devxp/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func LintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Run the set of linters defined in the .devxp config.",
		Args:  cobra.MaximumNArgs(0),
		Long: heredoc.Docf(`
			Run the set of linters defined in the .devxp config file. If no linters are defined,
			the command will do nothing.
		`, "`"),
		Example: heredoc.Doc(`
			// Create a new branch 'wip' and checkout to it:
			$ gh devxp lint
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			branchId := ""
			if len(args) > 0 {
				branchId = args[0]
			}

			return checkoutBranch(ctx, branchId)
		},
	}

	return cmd
}

func checkoutBranch(ctx context.Context, branchId string) error {
	// Does the branch exist?
	branchExists, err := branchExists(branchId)
	if err != nil {
		return err
	}

	if branchExists {
		log.Debugf("Branch '%s' already exists, checking out to it", branchId)
		out, err := utils.Exec("git", "checkout", branchId)
		if err != nil {
			return errors.Wrap(err, "Failed to checkout branch")
		}
		log.Info(strings.Trim(out, "\n"))
	} else {
		log.Debugf("Creating branch '%s' and checking out to it", branchId)
		out, err := utils.Exec("git", "checkout", "-b", branchId)
		if err != nil {
			return errors.Wrap(err, "Failed to create branch")
		}
		log.Info(strings.Trim(out, "\n"))
	}

	return nil
}

func branchExists(branchId string) (bool, error) {
	_, err := utils.Exec("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branchId)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 0 {
				return true, nil
			}
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
