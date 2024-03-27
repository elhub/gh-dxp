package diff

type Options struct {
	Confirm bool

	baseBranch string

	Reviewers []string
	Assignees []string
}

type PullRequest struct {
	Title string
	Body  string
}
