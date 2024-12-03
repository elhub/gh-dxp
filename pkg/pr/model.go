package pr

import (
	"github.com/charmbracelet/bubbles/table"
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
	branchID string
	Title    string
	Body     string
	isLinted bool
	isTested bool
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

type pullRequestInfo struct {
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

type pullRequestUI struct {
	table table.Model
}
