// pkg/status/status.go

package status

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/utils"
)

// Execute retrieves the status of the current repository.
func Execute(exe utils.Executor, opts *Options) error {
	var statusReport strings.Builder

	if optsIsEmpty(opts) {
		statusOptions := []string{"All", "Current repository", "PR status", "List branches", "Get all relevant Issues"}
		var selectedOption string
		survey.AskOne(&survey.Select{
			Message: "Choose the status type:",
			Options: statusOptions,
		}, &selectedOption)

		switch selectedOption {
		case "All":
			opts.All = true
		case "Current repository":
			opts.Repo = true
		case "PR status":
			opts.Pr = true
		case "List branches":
			opts.Branches = true
		case "Get all relevant Issues":
			opts.Issue = true
		}
	}

	if opts.All || opts.Repo {
		repo, err := exe.Command("git", "remote", "get-url", "origin")
		if err != nil {
			return err
		}
		statusReport.WriteString(fmt.Sprintf("Repository: %s\n", strings.TrimSpace(repo)))
	}

	if opts.All || opts.Pr {
		prStatus, err := exe.GH("pr", "status")
		if err != nil {
			return err
		}
		statusReport.WriteString(fmt.Sprintf("PR Status:\n%s\n", prStatus))
	}

	if opts.All || opts.Branches {
		branches, err := exe.Command("git", "branch", "-a")
		if err != nil {
			return err
		}
		statusReport.WriteString(fmt.Sprintf("Branches:\n%s\n", branches))
	}

	if opts.All || opts.Issue {
		assignedPRs, err := exe.GH("issue", "status")
		if err != nil {
			return err
		}
		statusReport.WriteString(fmt.Sprintf("Assigned PRs/Review Requests:\n%s\n", assignedPRs))
	}

	log.Info(statusReport.String())
	return nil
}

func optsIsEmpty(opts *Options) bool {
	return !opts.All && !opts.Repo && !opts.Pr && !opts.Branches && !opts.Issue
}
