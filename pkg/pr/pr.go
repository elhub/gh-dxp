// Package pr contains the logic for creating and updating pull requests.
package pr

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/branch"
	"github.com/elhub/gh-dxp/pkg/test"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/pkg/errors"
)

// Execute creates or updates a pull request, depending on its current state.
func Execute(exe utils.Executor, options *Options) error {
	// Get branchID
	currentBranch, errBranch := exe.Command("git", "branch", "--show-current")
	if errBranch != nil {
		return errBranch
	}
	branchID := strings.Trim(currentBranch, "\n")

	// Check if PR exists on branch
	prID, errCheck := CheckForExistingPR(exe, branchID)
	if errCheck != nil {
		return errCheck
	}

	if prID != "" {
		// If the PR exists, update it by pushing to the remote
		return update(exe, branchID, prID)
	}
	// If it doesn't exist, create a new PR
	return create(exe, options, branchID)
}

func create(exe utils.Executor, options *Options, branchID string) error {

	// Handle presence of untracked changes - ignore or abort
	untrackedChanges, err := getUntrackedChanges(exe)
	if err != nil {
		return err
	}

	if len(untrackedChanges) > 0 {
		if !options.AutoConfirm {
			res, err := askToConfirm(formatUntrackedFileChangesQuestion(untrackedChanges))
			if err != nil || !res {
				return errors.Wrap(err, "User aborted workflow")
			}
		}
	}

	//Handle presence of tracked changes - commit or abort
	trackedChanges, err := getTrackedChanges(exe)
	if err != nil {
		return err
	}

	if len(trackedChanges) > 0 {
		if !options.AutoConfirm {
			res, err := askToConfirm(formatTrackedFileChangesQuestion(trackedChanges))
			if err != nil || !res {
				return err
			}
			err = addAndCommitFiles(exe, trackedChanges)
			if err != nil {
				return err
			}

		}
	}

	//Run tests
	err = test.RunTest(exe)
	if err != nil {
		return err
	}

	// Push the current branch to git remote
	s := utils.StartSpinner("Pushing current branch to remote...", "Pushed working branch to remote.")
	currentBranch, err := exe.Command("git", "push", "--set-upstream", "origin", branchID)
	s.Stop()
	if err != nil {
		return err
	}
	log.Info("Current Branch:" + currentBranch + "\n")

	// Fetch the default branch
	baseBranch := options.baseBranch
	if baseBranch == "" {
		s = utils.StartSpinner("Fetching repository default branch...", "Fetched repository default branch")
		stdOut, errV := exe.GH("repo", "view", "--json", "defaultBranchRef", "--jq", ".defaultBranchRef.name")
		s.Stop()
		if errV != nil {
			return errors.Wrap(errV, "Failed to fetch default branch")
		}
		baseBranch = strings.Trim(stdOut.String(), "\n")
	}

	pr, err := createPR(exe, options, branchID, baseBranch)
	if err != nil {
		return err
	}

	s = utils.StartSpinner("Processing pull request...", "Pull request "+pr.Title+" created.")
	args := []string{"pr", "create", "--title", pr.Title, "--body", pr.Body, "--base", baseBranch}
	args = append(args, generatePRArgs(options)...)
	stdOut, err := exe.GH(args...)
	s.Stop()
	if err != nil {
		return errors.Wrap(err, "Failed to create pull request")
	}
	log.Info(strings.Trim(stdOut.String(), "\n"))

	return nil
}

func generatePRArgs(options *Options) []string {
	args := []string{}

	if len(options.Assignees) > 0 {
		args = append(args, "--assignee", strings.Join(options.Assignees, ","))
	}
	if len(options.Reviewers) > 0 {
		args = append(args, "--reviewer", strings.Join(options.Reviewers, ","))
	}

	return args
}

func update(exe utils.Executor, branchID string, prID string) error {

	//Run tests
	err := test.RunTest(exe)
	if err != nil {
		return err
	}

	// Push the current branch to the already existing git remote
	s := utils.StartSpinner("Updating Pull Request #"+prID+"...", "Pull Request #"+prID+" has been updated.")
	_, err = exe.Command("git", "push")
	s.Stop()
	if err != nil {
		return err
	}

	// Fetching this for info
	stdOut, err := exe.GH("pr", "list", "-H", branchID, "--json", "url", "--jq", ".[].url")
	if err != nil {
		return err
	}

	log.Info(strings.Trim(stdOut.String(), "\n") + "\n")

	return nil
}

// CheckForExistingPR checks if a PR already exists for the current branch.
func CheckForExistingPR(exe utils.Executor, branchID string) (string, error) {
	stdOut, err := exe.GH("pr", "list", "-H", branchID, "--json", "number", "--jq", ".[].number")

	if err != nil {
		log.Debug("Error: " + err.Error())
		return "", errors.New("Failed to find existing PR")
	}

	number := strings.Trim(stdOut.String(), "\n")

	return number, nil
}

// GetPRTitle gets the title of the current PR.
func GetPRTitle(exe utils.Executor) (string, error) {
	stdOut, err := exe.GH("pr", "view", "--json", "title", "--jq", ".title")

	if err != nil {
		return "", errors.New("Error getting PR title")
	}

	title := strings.Trim(stdOut.String(), "\n")

	return title, nil
}

func createPR(
	exe utils.Executor,
	options *Options,
	branchID string,
	mainID string,
) (PullRequest, error) {
	pr := PullRequest{}

	// Get the commit messages between the current branch and the main branch and put them in the PR body.
	commits, err := branch.GetCommitMessages(exe, mainID, branchID)
	if err != nil {
		return pr, err
	}

	// Get the title
	pr.Title = getDefaultTitle(commits)
	if !options.AutoConfirm {
		pr.Title, err = askForString("Title", pr.Title)
		if err != nil {
			return pr, err
		}
	}

	pr.Body, err = createBody(options, commits)
	if err != nil {
		return pr, err
	}

	logPullRequest(pr)

	return pr, nil
}

func createBody(options *Options, commits string) (string, error) {
	body := ""

	// Add a summary of the commits to the PR body
	commitLines := strings.Split(commits, "\n")
	if len(commitLines) > 1 {
		body = "Summary:\n"
		for _, line := range commitLines[1:] {
			body += "- " + line + "\n"
		}
	}

	bodySurvey := "No summary. Do you want to add one?"
	if body != "" {
		log.Info("Summary of commits:\n" + body)
		bodySurvey = "Do you want to change the summary?"
	}

	if !options.AutoConfirm {
		editBody, err := askToConfirm(bodySurvey)
		if err != nil {
			return "", err
		}

		if editBody {
			editedBody, errB := askForMultiline("Summary:\n")
			if errB != nil {
				return "", errB
			}
			body = "Summary:\n" + editedBody
		}
	}

	// TODO: What type of PR is this?
	// Multi-choice: Feature, Bug Fix, Documentation, Test, Refactor, Style, Build, Chore
	// Type should be set as a label.
	issueSection, err := issuesChanges(options)
	if err != nil {
		return "", err
	}

	body = addDocSection(body, issueSection)

	// TODO: Auto-Linting. If not auto-linted, ask why?
	// TODO: Auto-Testing. If not auto-tested, ask why?

	testSection, err := testingChanges(options)
	if err != nil {
		return "", err
	}

	body = addDocSection(body, testSection)

	docsSection, err := documentationChanges(options)
	if err != nil {
		return "", err
	}

	body = addDocSection(body, docsSection)

	return body, nil
}

func issuesChanges(options *Options) (string, error) {
	// Issue ID(s)
	// Optionally add the issue ID(s) to the PR body.
	body := ""
	if !options.AutoConfirm {
		issueIDString, errI := askForString("Issue IDs (seperate with commas):", "")
		if errI != nil {
			return "", errI
		}

		if issueIDString != "" {
			if body != "" {
				body += "\n"
			}

			issueIDs := strings.Split(issueIDString, ",")
			for i, id := range issueIDs {
				issueIDs[i] = strings.TrimSpace(id)
			}
			body += "Issue ID(s): " + strings.Join(issueIDs, ", ") + "\n"
		}
	}

	return body, nil
}

func testingChanges(options *Options) (string, error) {
	// Testing
	// TODO: Consider skipping this, if dealing with code where unit and integration
	// tests are not applicable. Replace with deploy test question?
	body := ""
	if !options.AutoConfirm {
		unitTestConfirm, err := askToConfirm("Did you add new unit tests?")
		if err != nil {
			return "", err
		}

		integrationTestConfirm, err := askToConfirm("Did you add new integration tests?")
		if err != nil {
			return "", err
		}

		testCommand, err := askForString("Test Command?", "")
		if err != nil {
			return "", err
		}
		body += "Testing:\n"
		body += "- [" + getCheckboxMark(unitTestConfirm) + "] Unit Tests\n"
		body += "- [" + getCheckboxMark(integrationTestConfirm) + "] Integration Tests\n"
		body += "- Test Command: " + testCommand + "\n"
	} else {
		body += "Testing:\n"
		body += "- [ ] Unit Tests\n"
		body += "- [ ] Integration Tests\n"
		body += "- Test Command:"
	}

	return body, nil
}

func documentationChanges(options *Options) (string, error) {
	// Documentation
	// Multi-choice: README.md, docs, storybook, no updates
	docOptions := []string{"No updates", "README.md", "docs", "storybook"}
	selectedDocs := []string{}
	if !options.AutoConfirm {
		err := survey.AskOne(&survey.MultiSelect{
			Message: "What documentation was updated?",
			Options: docOptions,
		}, &selectedDocs, survey.WithValidator(survey.Required))
		if err != nil {
			return "", err
		}
	} else {
		selectedDocs = append(selectedDocs, "No updates")
	}

	body := "\nDocumentation:\n"
	for _, doc := range selectedDocs {
		body += "- " + doc + "\n"
	}

	return body, nil
}

func addDocSection(body string, section string) string {
	if body != "" {
		body += "\n"
	}

	body += section
	return body
}

func askToConfirm(question string) (bool, error) {
	confirm := false
	err := survey.AskOne(&survey.Confirm{
		Message: question,
	}, &confirm, survey.WithValidator(survey.Required))
	if err != nil {
		return false, err
	}
	return confirm, nil
}

func askForString(question string, defaultAnswer string) (string, error) {
	var title string
	prompt := &survey.Input{
		Message: question,
		Default: defaultAnswer,
	}
	err := survey.AskOne(prompt, &title)
	if err != nil {
		return "", err
	}
	return title, nil
}

func askForMultiline(question string) (string, error) {
	lines := ""
	err := survey.AskOne(&survey.Multiline{
		Message: question,
	}, &lines, survey.WithValidator(survey.Required))
	if err != nil {
		return lines, err
	}
	return lines, nil
}

func getDefaultTitle(commits string) string {
	lines := strings.Split(commits, "\n")
	if len(lines) > 0 {
		return lines[0]
	}
	return ""
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

func filter(list []string, test func(string) bool) (ret []string) {
	for _, s := range list {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return ret
}

func getUntrackedChanges(exe utils.Executor) ([]string, error) {

	re := regexp.MustCompile(`^\?\?`)

	return getChanges(exe, re)
}

func getTrackedChanges(exe utils.Executor) ([]string, error) {

	re := regexp.MustCompile(`^([ADMRT]|\s)([ADMRT]|\s)\s`) //This regex is intended to catch all tracked changes except for unmerged conflicts

	return getChanges(exe, re)
}

func getChanges(exe utils.Executor, re *regexp.Regexp) ([]string, error) {
	changeString, err := exe.Command("git", "status", "--porcelain")
	if err != nil {
		return []string{}, err
	}

	changes := strings.Split(changeString, "\n")
	untrackedChanges := filter(changes, re.MatchString)

	//Remove the regex matched part of the string, leaving only the file name
	for i, s := range untrackedChanges {
		untrackedChanges[i] = re.ReplaceAllString(s, "")
	}

	return untrackedChanges, nil
}

func addAndCommitFiles(exe utils.Executor, files []string) error {
	commitMessage, err := askForString("Please enter a commit message", "")
	if err != nil {
		return err
	} else if len(commitMessage) == 0 {
		errors.New("Empty commit message not allowed")
	}

	_, err = exe.Command("git", "add", strings.Join(files, " "))
	if err != nil {
		return err
	}
	_, err = exe.Command("git", "commit", "-m", fmt.Sprintf(`"%s"`, commitMessage))
	if err != nil {
		return err
	}

	return nil
}

func formatUntrackedFileChangesQuestion(changes []string) string {

	return "You have untracked files locally \n\n" + strings.Join(changes, "\n") + "\n\nIgnore these files and continue?"
}

func formatTrackedFileChangesQuestion(changes []string) string {
	return "You have uncommitted files locally \n\n" + strings.Join(changes, "\n") + "\n\nDo you want to create a new commit with these changes?"
}
