package pr

// Options represents the options for the pr command.
type Options struct {
	AutoConfirm bool

	baseBranch string

	Reviewers []string
	Assignees []string
}

// PullRequest represents a pull request.
type PullRequest struct {
	Title string
	Body  string
}
