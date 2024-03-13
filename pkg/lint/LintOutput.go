package lint

type LinterOutput struct {
	Linter      string
	Path        string
	Line        int
	Column      int
	Description string
	Severity    string
	Source      string
}
