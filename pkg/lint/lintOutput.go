package lint

type LintOutput struct {
	Linter      string
	Path        string
	Line        int
	Character   int
	Code        string
	Description string
}
