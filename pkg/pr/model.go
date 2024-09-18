package pr

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

	baseBranch string

	Reviewers []string
	Assignees []string
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
