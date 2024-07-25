package cmd

import (
    // "fmt"
    "github.com/MakeNowJust/heredoc"
    // "github.com/AlecAivazis/survey/v2"
    "github.com/elhub/gh-dxp/pkg/status"
    "github.com/elhub/gh-dxp/pkg/utils"
    "github.com/spf13/cobra"
)

// StatusCmd creates a new cobra.Command for the status functionality.
func StatusCmd(exe utils.Executor) *cobra.Command {
    // var all, current, pr, branches, assigned bool
    opts := &status.Options{}

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
            # Interactive mode
            $ gh dxp status
        `),
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return status.Execute(exe, opts)
		},
	}

    fl := cmd.Flags()
    fl.BoolVarP(
        &opts.All,
        "all",
        "a",
        false,
        "Get all statuses",
    )
    fl.BoolVarP(
        &opts.Repo,
        "current",
        "c",
        false,
        "Get current repository target",
    )
    fl.BoolVarP(
        &opts.Pr,
        "pr",
        "p",
        false,
        "Get the PR status",
    )
    fl.BoolVarP(
        &opts.Branches,
        "branches",
        "b",
        false,
        "List all branches",
    )
    fl.BoolVarP(
        &opts.Issue,
        "issues",
        "i",
        false,
        "Get all relavant Issues",
    )



    //     RunE: func(_ *cobra.Command, _ []string) error {
    //         var statusType string

    //         if all {
    //             statusType.all = true
    //         } else if current {
    //             statusType = "Current repository"
    //         } else if pr {
    //             statusType = "PR status"
    //         } else if branches {
    //             statusType = "List branches"
    //         } else if assigned {
    //             statusType = "Assigned PRs/Review Requests"
    //         } else {
    //             // Interactive dialogue flow
    //             statusOptions := []string{"All", "Current repository", "PR status", "List branches", "Get all relavant Issues"}
    //             survey.AskOne(&survey.Select{
    //                 Message: "Choose the status type:",
    //                 Options: statusOptions,
    //             }, &statusType)
    //         }

    //         statusReport, err := status.Excecute(exe, statusType)
    //         if err != nil {
    //             return err
    //         }
    //         fmt.Println(statusReport)
    //         return nil
    //     },
    // }

    // Adding flags
    // cmd.Flags().BoolVarP(&all, "all", "a", false, "Get all statuses")
    // cmd.Flags().BoolVar(&pr, "pr", false, "Get the PR status")
    // cmd.Flags().BoolVar(&branches, "branches", false, "List all branches")
    // cmd.Flags().BoolVar(&assigned, "assigned", false, "Get all relavant Issues")

    return cmd
}