package lint

type lintOutput struct {
	linter      string
	path        string
	line        int
	character   int
	code        string
	description string
}
