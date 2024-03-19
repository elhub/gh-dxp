package lint

type LintOutput struct {
	Linter      string
	Path        string
	Line        int
	Column      int
	Description string
	Severity    string
	Source      string
}
