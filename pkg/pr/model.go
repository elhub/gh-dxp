package pr

// Options represents the options for the pr command.
type CreateOptions struct {
	TestRun bool
	NoLint  bool
	NoUnit  bool

	Branch        string
	CommitMessage string

	baseBranch string

	Reviewers []string
	Assignees []string
}

// Options represents the options for the merge command.
type MergeOptions struct {
	AutoConfirm bool
}

type UpdateOptions struct {
	TestRun bool
	NoLint  bool
	NoUnit  bool
}

// PullRequest represents a pull request.
type PullRequest struct {
	Title string
	Body  string
}
