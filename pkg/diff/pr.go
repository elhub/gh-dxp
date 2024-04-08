package diff

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/utils"
)

func createPR(
	exe utils.Executor,
	options *Options,
	branchID string,
	mainId string,
) (PullRequest, error) {
	log.Debug("Create the PR of the diff")
	log.Info("Params: " + branchID + " " + mainId)

	pr := PullRequest{}

	// Get the commit messages between the current branch and the main branch and put them in the PR body.
	// This is done by running `git log mainId..branchID --oneline` and adding the output to the PR body.
	// If the user has already pushed the branch to the remote, the PR body should include the commit messages.
	// If the user has not pushed the branch to the remote, the PR body should include a message that the branch has not been pushed yet.
	commits, err := exe.Command("git", "log", mainId+".."+branchID, "--oneline", "--pretty=format:%s")
	if err != nil {
		return pr, err
	}
	// Default PR title is first line of commit message
	lines := strings.Split(commits, "\n")
	defaultTitle := ""
	if len(lines) > 0 {
		defaultTitle = lines[0]
	}
	if !options.AutoConfirm {
		survey.AskOne(&survey.Input{
			Message: "Title:",
			Default: defaultTitle,
		}, &pr.Title, survey.WithValidator(survey.Required))
	} else {
		pr.Title = defaultTitle
	}

	// Add a summary of the commits to the PR body
	suggestedBody := ""
	if len(lines) > 1 {
		for i, line := range lines {
			lines[i] = "- " + line
		}
		suggestedBody = "Summary:\n"
		suggestedBody = strings.Join(lines[1:], "\n")
	}

	bodySurvey := ""
	if suggestedBody != "" {
		log.Info("Summary of commits:\n" + suggestedBody)
		bodySurvey = "Do you want to change the summary?"
	} else {
		bodySurvey = "No summary. Do you want to add one?"
	}

	editBody := false
	if !options.AutoConfirm {
		survey.AskOne(&survey.Confirm{
			Message: bodySurvey,
		}, &editBody, survey.WithValidator(survey.Required))
	}

	if editBody {
		editedBody := ""
		survey.AskOne(&survey.Multiline{
			Message: "Summary:",
		}, &editedBody, survey.WithValidator(survey.Required))
		pr.Body = "Summary:\n" + editedBody
	} else if suggestedBody != "" {
		pr.Body = suggestedBody
	}

	// TODO: What type of PR is this?
	// Multi-choice: Feature, Bug Fix, Documentation, Test, Refactor, Style, Build, Chore
	// Type should be set as a label.

	// Issue ID(s)
	// Optionally add the issue ID(s) to the PR body.
	issueIDString := ""
	if !options.AutoConfirm {
		survey.AskOne(&survey.Input{
			Message: "Issue ID(s):",
		}, &issueIDString)

		if pr.Body != "" {
			pr.Body += "\n"
		}
	}

	if issueIDString != "" {
		issueIDs := strings.Split(issueIDString, ",")
		pr.Body += "Issue ID(s): " + strings.Join(issueIDs, ", ") + "\n"
	}

	// Testing
	// TODO: Consider skipping this, if dealing with code where unit and integration
	// tests are not applicable. Replace with deploy test question?
	unitTestConfirm := false
	integrationTestConfirm := false
	testCommand := ""

	if !options.AutoConfirm {
		survey.AskOne(&survey.Confirm{
			Message: "Did you add new unit tests?",
		}, &unitTestConfirm, survey.WithValidator(survey.Required))

		survey.AskOne(&survey.Confirm{
			Message: "Did you add new integration tests?",
		}, &integrationTestConfirm, survey.WithValidator(survey.Required))

		survey.AskOne(&survey.Input{
			Message: "Test Command?",
		}, &testCommand, survey.WithValidator(survey.Required))
	}

	// TODO: Auto-Linting. If not auto-linted, ask why?
	// TODO: Auto-Testing. If not auto-tested, ask why?

	if pr.Body != "" {
		pr.Body += "\n"
	}

	pr.Body += "Testing:\n"
	pr.Body += "- [" + getCheckboxMark(unitTestConfirm) + "] Unit Tests\n"
	pr.Body += "- [" + getCheckboxMark(integrationTestConfirm) + "] Integration Tests\n"
	pr.Body += "- Test Command: " + testCommand + "\n"

	// Documentation
	// Multi-choice: README.md, docs, storybook, no updates
	docOptions := []string{"No updates", "README.md", "docs", "storybook"}
	selectedDocs := []string{}
	if !options.AutoConfirm {
		survey.AskOne(&survey.MultiSelect{
			Message: "What documentation was updated?",
			Options: docOptions,
		}, &selectedDocs, survey.WithValidator(survey.Required))
	} else {
		selectedDocs = append(selectedDocs, "No updates")
	}

	pr.Body += "\nDocumentation:\n"
	for _, doc := range selectedDocs {
		pr.Body += "- " + doc + "\n"
	}

	logPullRequest(pr)

	return pr, nil
}

func getCheckboxMark(confirm bool) string {
	if confirm {
		return "x"
	}
	return " "
}

func logPullRequest(pr PullRequest) {
	log.Info("Submitting the following pull request\n" + pr.Title + "\n\n" + pr.Body)
}
