// pkg/status/status.go

package status

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/elhub/gh-dxp/pkg/logger"
	"github.com/elhub/gh-dxp/pkg/ghutil"
)

// Execute retrieves the status of the current repository.
func Execute(exe ghutil.Executor, opts *Options) error {
	var statusReport strings.Builder

	if optsIsEmpty(opts) {
		if err := promptForOptions(opts); err != nil {
			return err
		}
	}

	if err := buildStatusReport(exe, opts, &statusReport); err != nil {
		return err
	}

	logger.Info(statusReport.String())
	return nil
}

func promptForOptions(opts *Options) error {
	statusOptions := []string{"All", "Current repository", "PR status", "List branches", "Get all relevant Issues"}
	var selectedOption string
	err := survey.AskOne(&survey.Select{
		Message: "Choose the status type:",
		Options: statusOptions,
	}, &selectedOption)
	if err != nil {
		return err
	}

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
	return nil
}

func buildStatusReport(exe ghutil.Executor, opts *Options, statusReport *strings.Builder) error {
	type section struct {
		enabled bool
		fn      func() error
	}

	sections := []section{
		{opts.All || opts.Repo, func() error { return appendRepo(exe, statusReport) }},
		{opts.All || opts.Pr, func() error { return appendPRStatus(exe, statusReport) }},
		{opts.All || opts.Branches, func() error { return appendBranches(exe, statusReport) }},
		{opts.All || opts.Issue, func() error { return appendIssues(exe, statusReport) }},
	}

	for _, s := range sections {
		if s.enabled {
			if err := s.fn(); err != nil {
				return err
			}
		}
	}

	return nil
}

func appendRepo(exe ghutil.Executor, statusReport *strings.Builder) error {
	repo, err := exe.Command("git", "remote", "get-url", "origin")
	if err != nil {
		return err
	}
	fmt.Fprintf(statusReport, "Repository: %s\n", strings.TrimSpace(repo))
	return nil
}

func appendPRStatus(exe ghutil.Executor, statusReport *strings.Builder) error {
	prStatus, err := exe.GH("pr", "status")
	if err != nil {
		return err
	}
	fmt.Fprintf(statusReport, "PR Status:\n%s\n", prStatus)
	return nil
}

func appendBranches(exe ghutil.Executor, statusReport *strings.Builder) error {
	branches, err := exe.Command("git", "branch", "-a")
	if err != nil {
		return err
	}
	fmt.Fprintf(statusReport, "Branches:\n%s\n", branches)
	return nil
}

func appendIssues(exe ghutil.Executor, statusReport *strings.Builder) error {
	assignedPRs, err := exe.GH("issue", "status")
	if err != nil {
		return err
	}
	fmt.Fprintf(statusReport, "Assigned PRs/Review Requests:\n%s\n", assignedPRs)
	return nil
}

func optsIsEmpty(opts *Options) bool {
	return !opts.All && !opts.Repo && !opts.Pr && !opts.Branches && !opts.Issue
}
