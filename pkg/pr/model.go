// Package pr contains the functions and types for the pull request command.
package pr

import (
	"charm.land/bubbles/v2/table"
)

// Options represents the options for the pr command.
type Options struct {
	TestRun bool
	NoLint  bool
	NoUnit  bool

	CommitMessage string
}

// CreateOptions represents the options for the pr create command.
type CreateOptions struct {
	TestRun bool
	NoLint  bool
	NoUnit  bool
	Draft   bool

	Branch        string
	CommitMessage string
	Issues        string

	baseBranch string

	Reviewers []string
	Assignees []string
}

// ListOptions represents the options for the pr list command.
type ListOptions struct {
	TestRun         bool
	Mine            bool
	ReviewRequested bool
}

// MergeOptions represents the options for the pr merge command.
type MergeOptions struct {
	AutoConfirm bool
}

// UpdateOptions represents the options for the pr update command.
type UpdateOptions struct {
	TestRun bool
	NoLint  bool
	NoUnit  bool

	CommitMessage string
}

// PullRequest represents a pull request.
type PullRequest struct {
	branchID 		string
	targetBranch    string
	Title    		string
	Body     		string
	isLinted 		bool
	isTested 		bool
	label 			string
}

// The following structs are used to unmarshal the JSON responses from the GitHub API.
type searchResult struct {
	Number     int              `json:"number"`
	Repository searchRepository `json:"repository"`
}

type searchRepository struct {
	Name          string `json:"name"`
	NameWithOwner string `json:"nameWithOwner"`
}

// PullRequestInfo represents detailed information about a pull request.
type PullRequestInfo struct {
	Additions      int          `json:"additions"`
	Author         prAuthor     `json:"author"`
	CreatedAt      string       `json:"createdAt"`
	Deletions      int          `json:"deletions"`
	HeadRepository prRepository `json:"headRepository"`
	Number         int          `json:"number"`
	ReviewDecision string       `json:"reviewDecision"`
	Title          string       `json:"title"`
}

type prRepository struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type prAuthor struct {
	ID    string `json:"id"`
	IsBot bool   `json:"is_bot"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

// PullRequestUI represents the UI model for displaying pull requests.
type PullRequestUI struct {
	table table.Model
}

// PullRequestLabel represents the model for the pull request label.
type PullRequestLabel struct {
	Name  		string
	Color 		string
	Description string
}

// PullRequestLabels defines the labels that can be applied to pull requests.
var PullRequestLabels = []PullRequestLabel {
	{
		Name: "Bugfix",
		Color: "#d93f0b",
		Description: "A bugfix PR is a pull request that fixes a bug in the codebase.",
	},
	{
		Name: "Build",
		Color: "#196020",
		Description: "A build PR is a pull request that updates the build system or dependencies.",
	},
	{
		Name: "Chore",
		Color: "#5319e7",
		Description: "A chore PR is a pull request that performs routine tasks or maintenance.",
	},
	{
		Name: "Documentation",
		Color: "#dda5fc",
		Description: "A documentation PR is a pull request that updates the documentation.",
	},
	{
		Name: "Feature",
		Color: "#0075ca",
		Description: "A feature PR is a pull request that adds a new feature to the codebase.",
	},
	{
		Name: "Refactor",
		Color: "#e99695",
		Description: "A refactor PR is a pull request that restructures existing code without changing its behavior.",
	},
	{
		Name: "Style",
		Color: "#17f7e9",
		Description: "A style PR is a pull request that changes the code style or formatting.",
	},
	{
		Name: "Test",
		Color: "#fbca04",
		Description: "A test PR is a pull request that adds or updates tests.",
	},
}
