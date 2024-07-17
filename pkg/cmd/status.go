package cmd

import (
    "fmt"
	"github.com/MakeNowJust/heredoc"
    "github.com/AlecAivazis/survey/v2"
    "github.com/elhub/gh-dxp/pkg/status"
    "github.com/elhub/gh-dxp/pkg/utils"
    "github.com/spf13/cobra"
)

// NewStatusCmd creates a new cobra.Command for the status functionality.
func StatusCmd(exec utils.Executor) *cobra.Command {
    var all, current, pr, branches, assigned bool

    cmd := &cobra.Command{
        Use:   "status",
        Short: "Get the status of the current repository",
        Long: heredoc.Docf(`
            Get the status of the current repository or specific aspects related to it. This command allows you to:

            * View a comprehensive status report covering all aspects (All)
            * Check the status of the current repository
            * View the status of pull requests (PRs)
            * List all branches within the repository
            * View assigned PRs/Review Requests

            This command supports both interactive mode and non-interactive mode via flags for quick access to specific information.
        `, "`"),
        Example: heredoc.Doc(`
            $ gh dxp status --all
            $ gh dxp status --current
            $ gh dxp status --pr
            $ gh dxp status --branches
            $ gh dxp status --assigned

            # Interactive mode
            $ gh dxp status
        `),
        RunE: func(cmd *cobra.Command, args []string) error {
            statusChecker := status.NewStatus(exec)
            var statusType string

            if all {
                statusType = "All"
            } else if current {
                statusType = "current repository"
            } else if pr {
                statusType = "PR status"
            } else if branches {
                statusType = "List branches"
            } else if assigned {
                statusType = "assigned PRs/Review Requests"
            } else {
                // Interactive dialogue flow
                statusOptions := []string{"All", "current repository", "PR status", "List branches", "assigned PRs/Review Requests"}
                survey.AskOne(&survey.Select{
                    Message: "Choose the status type:",
                    Options: statusOptions,
                }, &statusType)
            }

            statusReport, err := statusChecker.GetStatus(statusType)
            if err != nil {
                return err
            }
            fmt.Println(statusReport)
            return nil
        },
    }

    // Adding flags
    cmd.Flags().BoolVarP(&all, "all", "a", false, "Get all statuses")
    cmd.Flags().BoolVar(&current, "current", false, "Get the current repository status")
    cmd.Flags().BoolVar(&pr, "pr", false, "Get the PR status")
    cmd.Flags().BoolVar(&branches, "branches", false, "List all branches")
    cmd.Flags().BoolVar(&assigned, "assigned", false, "Get assigned PRs/Review Requests")

    return cmd
}