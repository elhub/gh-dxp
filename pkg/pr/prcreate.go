// Package pr provides utilities for managing pull requests in gh-dxp.
package pr

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/elhub/gh-dxp/pkg/branch"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/ghutil"
	"github.com/elhub/gh-dxp/pkg/jira"
	"github.com/elhub/gh-dxp/pkg/logger"
	"github.com/pkg/errors"
)

// jiraKeyPattern matches Jira issue keys like TDX-123 and EDIEL-456.
var jiraKeyPattern = regexp.MustCompile(`\b[A-Z][A-Z0-9]+-\d+\b`)

// ExtractJiraIDs parses a branch name and returns any Jira issue keys found.
func ExtractJiraIDs(branchName string) []string {
	return jiraKeyPattern.FindAllString(branchName, -1)
}

// CreateTemporaryBranch creates a new temporary branch from the current base branch and checks it out. It updates the pr struct with the new branch name.
func CreateTemporaryBranch(exe ghutil.Executor, options *CreateOptions, pr *PullRequest) error {
	newBranchName, err := getNewBranchName(options)
	if err != nil {
		return err
	}

	branchExists, err := branch.Exists(exe, newBranchName)
	if err != nil {
		return err
	}
	if branchExists {
		return errors.New("Branch already exists. Please delete it or specify another one")
	}
	_, err = exe.Command("git", "checkout", "-b", newBranchName)
	if err != nil {
		return err
	}
	pr.branchID = newBranchName
	return nil
}

// ExecuteCreate creates or updates a pull request, depending on its current state.
func ExecuteCreate(exe ghutil.Executor, settings *config.Settings, options *CreateOptions) error {
	pr := PullRequest{}
	// Get branchID
	currentBranch, errBranch := exe.Command("git", "branch", "--show-current")
	if errBranch != nil {
		return errBranch
	}
	pr.branchID = strings.Trim(currentBranch, "\n")

	baseBranch, err := setBaseBranch(exe, options)
	if err != nil {
		return err
	}
	if baseBranch == "" {
		return errors.New("Base branch cannot be empty")
	}
	pr.targetBranch = baseBranch

	// If we're currently in the base branch, we need to make a new temporary branch to contain the diff
	if pr.branchID == pr.targetBranch {
		err := CreateTemporaryBranch(exe, options, &pr)
		if err != nil {
			return err
		}
	}

	// Check if PR exists on branch
	prID, errCheck := CheckForExistingPR(exe, pr.branchID)
	if errCheck != nil {
		return errCheck
	}

	prOpts := &Options{
		TestRun:       options.TestRun,
		NoLint:        options.NoLint,
		NoUnit:        options.NoUnit,
		CommitMessage: options.CommitMessage,
	}

	pr, err = performPreCreateOperations(exe, settings, pr, prOpts)
	if err != nil {
		return err
	}

	if prID != "" {
		// If the PR exists, update it by pushing to the remote
		return update(exe, pr.branchID, prID)
	}

	// If it doesn't exist, create a new PR
	return create(exe, options, settings, pr)
}

func create(exe ghutil.Executor, options *CreateOptions, settings *config.Settings, pr PullRequest) error {
	// Push the current branch to git remote
	s := ghutil.StartSpinner("Pushing current branch to remote...", "Pushed working branch to remote.")
	currentBranch, err := exe.Command("git", "push", "--set-upstream", "origin", pr.branchID)
	if err != nil {
		ghutil.RemoveFinalMsg(s)
		return err
	}
	s.Stop()
	logger.Info("Current Branch:" + currentBranch + "\n")
	newPR, err := createPR(exe, options, settings, pr, options.baseBranch)
	if err != nil {
		return err
	}

	err = ensureLabelExistsInRepository(exe, pr.label)
	if err != nil {
		return err
	}

	s = ghutil.StartSpinner("Processing pull request...", "Pull request "+newPR.Title+" created.")
	args := []string{"pr", "create", "--title", newPR.Title, "--body", newPR.Body, "--base", options.baseBranch, "--label", pr.label}
	args = append(args, generatePRArgs(options)...)
	stdOut, err := exe.GH(args...)
	if err != nil {
		ghutil.RemoveFinalMsg(s)
		return errors.Wrap(err, "Failed to create pull request")
	}
	s.Stop()
	logger.Info(strings.Trim(stdOut, "\n"))

	return nil
}

func ensureLabelExistsInRepository(exe ghutil.Executor, labelName string) error {
	stdOut, err := exe.GH("label", "list", "--limit", "1000", "--json", "name", "--jq", ".[].name")
	if err != nil {
		return err
	}
	if strings.Contains(stdOut, labelName) {
		return nil
	}
	var label PullRequestLabel
	found := false
	for _, l := range PullRequestLabels {
		if l.Name == labelName {
			label = l
			found = true
			break
		}
	}
	if !found {
		return errors.Errorf("unrecognized pull request label: %s", labelName)
	}
	logger.Info("Label " + labelName + " does not exist in repository. Creating label \"" + label.Name + "\"...")
	_, err = exe.GH("label", "create", label.Name, "--color", label.Color, "--description", label.Description)
	if err != nil {
		return err
	}
	return nil
}

func generatePRArgs(options *CreateOptions) []string {
	args := []string{}

	if len(options.Assignees) > 0 {
		args = append(args, "--assignee", strings.Join(options.Assignees, ","))
	}
	if len(options.Reviewers) > 0 {
		args = append(args, "--reviewer", strings.Join(options.Reviewers, ","))
	}
	if options.Draft {
		args = append(args, "--draft")
	}

	return args
}

func update(exe ghutil.Executor, branchID string, prID string) error {
	// Push the current branch to the already existing git remote
	s := ghutil.StartSpinner("Updating Pull Request #"+prID+"...", "Pull Request #"+prID+" has been updated.")
	_, err := exe.Command("git", "push")
	if err != nil {
		ghutil.RemoveFinalMsg(s)
		return err
	}
	s.Stop()

	// Fetching this for info
	stdOut, err := exe.GH("pr", "list", "-H", branchID, "--json", "url", "--jq", ".[].url")
	if err != nil {
		return err
	}

	logger.Info(strings.Trim(stdOut, "\n") + "\n")

	return nil
}

func createPR(
	exe ghutil.Executor,
	options *CreateOptions,
	settings *config.Settings,
	pr PullRequest,
	mainID string,
) (PullRequest, error) {
	// Get the commit messages between the current branch and the main branch and put them in the PR body.
	commits, err := branch.GetCommitMessages(exe, mainID, pr.branchID)
	if err != nil {
		return pr, err
	}

	// Get the title
	pr.Title = getDefaultTitle(commits)
	if !options.TestRun {
		pr.Title, err = ghutil.AskForString("Title", pr.Title)
		if err != nil {
			return pr, err
		}
	}

	pr.Body, err = createBody(exe, pr, options, settings, commits)
	if err != nil {
		return pr, err
	}

	logPullRequest(pr)

	return pr, nil
}

func createBody(exe ghutil.Executor, pr PullRequest, options *CreateOptions, settings *config.Settings, commits string) (string, error) {
	body := ""

	// Add a summary of the commits to the PR body
	commitLines := strings.Split(commits, "\n")
	commitSummary := ""
	if len(commitLines) > 1 {
		var commitSummarySb185 strings.Builder
		for _, line := range commitLines[1:] {
			commitSummarySb185.WriteString("* " + line + "\n")
		}
		commitSummary += commitSummarySb185.String()
	}

	bodySurvey := "No description. Do you want to add one?"
	if body != "" {
		logger.Info("## Description\n\n" + commitSummary)
		bodySurvey = "Do you want to change the description?"
	}

	if !options.TestRun {
		var err error
		body, err = promptForDescription(body, commitSummary, bodySurvey)
		if err != nil {
			return "", err
		}
	}

	issueSection, err := issuesChanges(options, settings, pr.branchID, commits, pr.Title, body)
	if err != nil {
		return "", err
	}

	body = addDocSection(body, issueSection)

	// CheckList
	body = addDocSection(body, "## 📋 Checklist\n")

	body = addDocSection(body, docIsLintedLine(pr, options))
	body = addDocSection(body, docIsTestedLine(pr, options))

	// New tests checkmark
	testSection, err := testingChanges(options)
	if err != nil {
		return "", err
	}
	body = addDocSection(body, testSection)

	docsSection, err := documentationChanges(exe)
	if err != nil {
		return "", err
	}
	body = addDocSection(body, docsSection)

	// POSIX - always end with \n
	// Append a newline to the end of the body if it does not have one
	if !strings.HasSuffix(body, "\n") {
		body += "\n"
	}

	return body, nil
}

func promptForDescription(body, commitSummary, bodySurvey string) (string, error) {
	editBody, err := ghutil.AskToConfirm(bodySurvey)
	if err != nil {
		return "", err
	}
	if editBody {
		editedBody, errB := ghutil.AskForMultiline("Description:\n")
		if errB != nil {
			return "", errB
		}
		return "## 📝 Description\n\n" + editedBody + "\n", nil
	}
	if commitSummary != "" {
		return "## 📝 Description\n\n" + commitSummary, nil
	}
	return body, nil
}

func docIsTestedLine(pr PullRequest, options *CreateOptions) string {
	switch {
	case pr.isTested:
		return "* ✅ Unit tests passed on local machine."
	case options.NoUnit:
		return "* ⛔ **This PR has not been unit tested! The --notest option was used.**"
	default:
		return "* ⚠️ **No tests could be run for this PR.**"
	}
}

func docIsLintedLine(pr PullRequest, options *CreateOptions) string {
	switch {
	case pr.isLinted:
		return "* ✅ Lint checks passed on local machine."
	case options.NoLint:
		return "* ⛔ **This PR has not been linted! The --nolint option was used.**"
	default:
		return "* ⛔ **This PR has not been linted! Unspecified lint error!** ⚠️"
	}
}

func issuesChanges(options *CreateOptions, settings *config.Settings, branchName, commits, title, description string) (string, error) {
	// Issue ID(s)
	// Optionally add the issue ID(s) to the PR body.
	body := ""
	var issueIDString string
	autoDetected := false
	detectedIDs := ExtractJiraIDs(branchName)
	if !options.TestRun && options.Issues == "" {
		if suggestions, err := jira.SearchIssues(context.Background(), settings.JiraURL, settings.JiraEmail, jira.SearchText{
			CommitMessage: commits,
			Title:         title,
			Description:   description,
		}); err != nil {
			switch err {
			case jira.ErrJiraNotConfigured:
				logger.Info("💡 Tip: Add your Jira API token to ~/.jira_token to get automatic ticket suggestions.\n   Create an empty ~/.jira_token to hide this message.")
			case jira.ErrJiraDisabled:
				// user opted out — stay silent
			default:
				logger.Warn("Unable to fetch Jira suggestions: " + err.Error())
			}
		} else if len(suggestions) > 0 {
			logger.Info(formatJiraSuggestions(suggestions))
			seen := map[string]bool{}
			for _, id := range detectedIDs {
				seen[id] = true
			}
			for i, s := range suggestions {
				if i == 3 {
					break
				}
				if !seen[s.Key] {
					detectedIDs = append(detectedIDs, s.Key)
				}
			}
		}

		userIssueString, errI := ghutil.AskForString(
			"Issue IDs (separate with commas):",
			strings.Join(detectedIDs, ", "),
		)
		if errI != nil {
			return "", errI
		}
		issueIDString = userIssueString
	} else if options.Issues == "" {
		autoDetected = true
		issueIDString = strings.Join(detectedIDs, ", ")
	} else {
		issueIDString = options.Issues
	}
	if issueIDString != "" {
		if body != "" {
			body += "\n"
		}

		issueIDs := strings.Split(issueIDString, ",")
		for i := range issueIDs {
			issueIDs[i] = strings.TrimSpace(issueIDs[i])
			if !autoDetected {
				issueIDs[i] = fmt.Sprintf("[%s](%s/%s)", issueIDs[i], settings.JiraURL, issueIDs[i])
			}
		}
		body += "## 🔗 Issue ID(s): " + strings.Join(issueIDs, ", ") + "\n"
	}

	return body, nil
}

func formatJiraSuggestions(issues []jira.SearchIssue) string {
	var suggestionLines []string
	for i, issue := range issues {
		if i == 5 {
			break
		}
		suggestionLines = append(suggestionLines, fmt.Sprintf("%d. %s - %s", i+1, issue.Key, issue.Fields.Summary))
	}
	if len(suggestionLines) == 0 {
		return ""
	}
	return "Suggested issues:\n" + strings.Join(suggestionLines, "\n")
}

func testingChanges(options *CreateOptions) (string, error) {
	if !options.TestRun {
		newTestConfirm, err := ghutil.AskToConfirm("Did you add new tests?")
		if err != nil {
			return "", err
		}

		if newTestConfirm {
			return "* ✅ This PR adds new tests.", nil
		}
	}

	return "", nil
}

func documentationChanges(exe ghutil.Executor) (string, error) {
	changedFiles, err := ghutil.GetChangedFiles(exe)
	if err != nil {
		return "", err
	}

	readmewasUpdated := ghutil.CheckFilesUpdated(changedFiles, []string{"README.md$"})
	docsWereUpdated := ghutil.CheckFilesUpdated(changedFiles, []string{"/docs/"})

	selectedDocs := []string{}
	if readmewasUpdated {
		selectedDocs = append(selectedDocs, "README")
	}
	if docsWereUpdated {
		selectedDocs = append(selectedDocs, "System Documentation")
	}

	body := ""
	if len(selectedDocs) > 0 {
		body += "* ✅ Documentation Updates: " + strings.Join(selectedDocs, ", ")
	}

	return body, nil
}

func addDocSection(body, section string) string {
	if (section == "") || (section == "\n") {
		return body
	}

	if body != "" {
		body += "\n"
	}

	body += section
	return body
}

func getDefaultTitle(commits string) string {
	lines := strings.Split(commits, "\n")
	if len(lines) > 0 {
		return lines[0]
	}
	return ""
}

func logPullRequest(pr PullRequest) {
	logger.Info("Submitting the following pull request\n" + pr.Title + "\n\n" + pr.Body)
}

func formatUntrackedFileChangesQuestion(changes []string) string {
	return "You have untracked files locally \n\n" + strings.Join(changes, "\n") + "\n\nIgnore these files and continue?"
}

func formatTrackedFileChangesQuestion(changes []string) string {
	return "You have uncommitted files locally \n\n" + strings.Join(changes, "\n") + "\n\nDo you want to create a new commit with these changes?"
}

// If the baseBranch option is not set, set it to the base branch of the remote.
func setBaseBranch(exe ghutil.Executor, options *CreateOptions) (string, error) {
	// Fetch the default branch
	baseBranch := options.baseBranch
	if baseBranch == "" {
		s := ghutil.StartSpinner("Fetching repository default branch...", "Fetched repository default branch")
		stdOut, errV := exe.GH("repo", "view", "--json", "defaultBranchRef", "--jq", ".defaultBranchRef.name")
		if errV != nil {
			ghutil.RemoveFinalMsg(s)
			return "", errors.Wrap(errV, "Failed to fetch default branch")
		}
		s.Stop()
		baseBranch = strings.Trim(stdOut, "\n")
		options.baseBranch = baseBranch
	}
	return baseBranch, nil
}

func getNewBranchName(options *CreateOptions) (string, error) {
	var newBranchName = "branch1"

	if options.Branch != "" {
		return options.Branch, nil
	}

	if !options.TestRun {
		inputBranchName, err := ghutil.AskForString("You are currently on the base branch. Please specify a temporary branch name: ", "")
		if err != nil {
			return "", err
		}
		if inputBranchName == "" {
			return "", errors.New("Branch name cannot be empty")
		}
		newBranchName = inputBranchName
	}
	return newBranchName, nil
}
