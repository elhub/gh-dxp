package pr

import (
	"path/filepath"
	"strings"

	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/branch"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/lint"
	"github.com/elhub/gh-dxp/pkg/test"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/pkg/errors"
)

// ExecuteCreate creates or updates a pull request, depending on its current state.
func ExecuteCreate(exe utils.Executor, settings *config.Settings, options *CreateOptions) error {
	// Get branchID
	currentBranch, errBranch := exe.Command("git", "branch", "--show-current")
	if errBranch != nil {
		return errBranch
	}
	branchID := strings.Trim(currentBranch, "\n")

	baseBranch, err := setBaseBranch(exe, options)
	if err != nil {
		return err
	}

	// If we're currently in the base branch, we need to make a new temporary branch to contain the diff
	if branchID == baseBranch {
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
		branchID = newBranchName
	} else {
		if options.Branch != "" && options.Branch != branchID {
			log.Info("Branch option was specified, but we are not currently on the default branch. Proceeding with branch " + branchID)
		}
	}

	// Check if PR exists on branch
	prID, errCheck := CheckForExistingPR(exe, branchID)
	if errCheck != nil {
		return errCheck
	}

	err = performPreCommitOperations(exe, settings, options)
	if err != nil {
		return err
	}

	if prID != "" {
		// If the PR exists, update it by pushing to the remote
		return update(exe, branchID, prID)
	}
	// If it doesn't exist, create a new PR
	return create(exe, options, branchID)
}

func performPreCommitOperations(exe utils.Executor, settings *config.Settings, options *CreateOptions) error {
	// Handle uncommitted changes
	filesToCommit, err := handleUncommittedChanges(exe, options)
	if err != nil {
		return err
	}

	// Run tests
	if !options.NoUnit {
		err = test.RunTest(exe)
		if err != nil {
			return err
		}
	}

	// Run lint
	if !options.NoLint {
		err = lint.Run(exe, settings, &lint.Options{})
		if err != nil {
			return err
		}
	}

	if len(filesToCommit) > 0 {
		err = addAndCommitFiles(exe, filesToCommit, options)
		if err != nil {
			return err
		}
	}

	return nil
}

func create(exe utils.Executor, options *CreateOptions, branchID string) error {
	// Push the current branch to git remote
	s := utils.StartSpinner("Pushing current branch to remote...", "Pushed working branch to remote.")
	currentBranch, err := exe.Command("git", "push", "--set-upstream", "origin", branchID)
	s.Stop()
	if err != nil {
		return err
	}
	log.Info("Current Branch:" + currentBranch + "\n")
	pr, err := createPR(exe, options, branchID, options.baseBranch)
	if err != nil {
		return err
	}

	s = utils.StartSpinner("Processing pull request...", "Pull request "+pr.Title+" created.")
	args := []string{"pr", "create", "--title", pr.Title, "--body", pr.Body, "--base", options.baseBranch}
	args = append(args, generatePRArgs(options)...)
	stdOut, err := exe.GH(args...)
	s.Stop()
	if err != nil {
		return errors.Wrap(err, "Failed to create pull request")
	}
	log.Info(strings.Trim(stdOut.String(), "\n"))

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

func update(exe utils.Executor, branchID string, prID string) error {
	// Push the current branch to the already existing git remote
	s := utils.StartSpinner("Updating Pull Request #"+prID+"...", "Pull Request #"+prID+" has been updated.")
	_, err := exe.Command("git", "push")
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

func createPR(
	exe utils.Executor,
	options *CreateOptions,
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
	if !options.TestRun {
		pr.Title, err = utils.AskForString("Title", pr.Title)
		if err != nil {
			return pr, err
		}
	}

	pr.Body, err = createBody(exe, options, commits)
	if err != nil {
		return pr, err
	}

	logPullRequest(pr)

	return pr, nil
}

func createBody(exe utils.Executor, options *CreateOptions, commits string) (string, error) {
	body := ""

	// Add a summary of the commits to the PR body
	commitLines := strings.Split(commits, "\n")
	commitSummary := ""
	if len(commitLines) > 1 {
		for _, line := range commitLines[1:] {
			commitSummary += "* " + line + "\n"
		}
	}

	bodySurvey := "No description. Do you want to add one?"
	if body != "" {
		log.Info("## Description\n\n" + commitSummary)
		bodySurvey = "Do you want to change the description?"
	}

	if !options.TestRun {
		editBody, err := utils.AskToConfirm(bodySurvey)
		if err != nil {
			return "", err
		}

		if editBody {
			editedBody, errB := utils.AskForMultiline("Description:\n")
			if errB != nil {
				return "", errB
			}
			body = "## 📝 Description\n\n" + editedBody + "\n"
		} else {
			body = "## 📝 Description\n\n" + commitSummary
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

	// CheckList
	body = addDocSection(body, "## 📋 Checklist\n")

	if options.NoLint {
		body = addDocSection(body, "* ⛔ **This PR has not been linted!** ⚠️")
		// TODO: Auto-Testing. If not auto-tested, ask why?
	} else {
		body = addDocSection(body, "* ✅ Lint checks passed on local machine.")
	}
	if options.NoUnit {
		body = addDocSection(body, "* ⛔ **This PR has not been unit tested!** ⚠️")
		// TODO: Auto-Linting. If not auto-linted, ask why?
	} else {
		body = addDocSection(body, "* ✅ Unit tests passed on local machine.")
	}

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

func issuesChanges(options *CreateOptions) (string, error) {
	// Issue ID(s)
	// Optionally add the issue ID(s) to the PR body.
	body := ""
	if !options.TestRun {
		issueIDString, errI := utils.AskForString("Issue IDs (seperate with commas):", "")
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
			body += "## 🔗 Issue ID(s): " + strings.Join(issueIDs, ", ") + "\n"
		}
	}

	return body, nil
}

func testingChanges(options *CreateOptions) (string, error) {
	if !options.TestRun {
		newTestConfirm, err := utils.AskToConfirm("Did you add new tests?")
		if err != nil {
			return "", err
		}

		if newTestConfirm {
			return "* ✅ This PR adds news tests.", nil
		}

	}

	return "", nil
}

func handleUncommittedChanges(exe utils.Executor, options *CreateOptions) ([]string, error) {
	// Handle presence of untracked changes - ignore or abort
	untrackedChanges, err := utils.GetUntrackedChanges(exe)
	if err != nil {
		return []string{}, err
	}

	if len(untrackedChanges) > 0 && !options.TestRun {
		res, err := utils.AskToConfirm(formatUntrackedFileChangesQuestion(untrackedChanges))
		if err != nil {
			return []string{}, err
		}
		if !res {
			return []string{}, errors.New("User aborted workflow")
		}
	}

	// Handle presence of tracked changes - commit or abort
	trackedChanges, err := utils.GetTrackedChanges(exe)
	if err != nil {
		return []string{}, err
	}

	if len(trackedChanges) > 0 && !options.TestRun && options.CommitMessage == "" {
		res, err := utils.AskToConfirm(formatTrackedFileChangesQuestion(trackedChanges))
		if err != nil {
			return []string{}, err
		}

		if !res {
			return []string{}, errors.New("User aborted workflow")
		}
	}
	return trackedChanges, nil
}

func documentationChanges(exe utils.Executor) (string, error) {
	changedFiles, err := utils.GetChangedFiles(exe)
	if err != nil {
		return "", err
	}

	readmewasUpdated := utils.CheckFilesUpdated(changedFiles, []string{"README.md$"})
	docsWereUpdated := utils.CheckFilesUpdated(changedFiles, []string{"/docs/"})

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

func addDocSection(body string, section string) string {
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

func getCheckboxMark(confirm bool) string {
	if confirm {
		return "x"
	}
	return " "
}

func logPullRequest(pr PullRequest) {
	log.Info("Submitting the following pull request\n" + pr.Title + "\n\n" + pr.Body)
}

func addAndCommitFiles(exe utils.Executor, files []string, options *CreateOptions) error {
	var commitMessage string
	var err error

	if options.CommitMessage != "" {
		commitMessage = options.CommitMessage
	} else {

		if !options.TestRun {
			commitMessage, err = utils.AskForString("Please enter a commit message: ", "")
			if err != nil {
				return err
			} else if len(commitMessage) == 0 {
				return errors.New("Empty commit message not allowed")
			}
		} else {
			commitMessage = "default commit message"
		}
	}
	// Get git root directory and add to files to get fully qualified paths
	root, err := utils.GetGitRootDirectory(exe)
	if err != nil {
		return err
	}

	var fullPaths []string
	for _, filePath := range files {
		fullPaths = append(fullPaths, filepath.Join(root, filePath))
	}

	addCommandArgs := append([]string{"add"}, fullPaths...)

	_, err = exe.Command("git", addCommandArgs...)
	if err != nil {
		return err
	}

	// Commit files
	_, err = exe.Command("git", "commit", "-m", commitMessage)
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

// If the baseBranch option is not set, set it to the base branch of the remote.
func setBaseBranch(exe utils.Executor, options *CreateOptions) (string, error) {
	// Fetch the default branch
	baseBranch := options.baseBranch
	if baseBranch == "" {
		s := utils.StartSpinner("Fetching repository default branch...", "Fetched repository default branch")
		stdOut, errV := exe.GH("repo", "view", "--json", "defaultBranchRef", "--jq", ".defaultBranchRef.name")
		s.Stop()
		if errV != nil {
			return "", errors.Wrap(errV, "Failed to fetch default branch")
		}
		baseBranch = strings.Trim(stdOut.String(), "\n")
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
		inputBranchName, err := utils.AskForString("You are currently on the base branch. Please specify a temporary branch name: ", "")
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
