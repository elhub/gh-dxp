package pr

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/caarlos0/log"
)

type PullRequest struct {
	Title  string
	Body   string
	Labels []string
}

func createPR(
	branchID string,
	confirm bool,
) (*PullRequest, error) {
	log.Debug("Create the PR")

	pr := PullRequest{}
	pr.Title = branchID

	// What type of PR is this?
	// Multi-choice: Feature, Bug Fix, Documentation, Test, Refactor, Style, Build, Chore
	// Type is set as a label.

	// Issue ID(s)
	// Add the issue ID(s) to the PR body.
	issueIDs := []string{}
	survey.AskOne(&survey.Input{
		Message: "Issue ID(s)",
	}, &issueIDs, survey.WithValidator(survey.Required))
	pr.Body += "Issue ID(s): " + strings.Join(issueIDs, ", ") + "\n"

	// Test Plan
	// Test plan is described in the PR body.
	// If auto-tested, suggest the auto-test as the test plan.

	// Documentation
	// Multi-choice: README.md, docs, storybook, none needed

	// Linting
	// If not auto-linted, ask why?

	// Testing
	// If not auto-tested, ask why?

	// Confirm body or edit

	return pr, nil
}
