package lint

type Linter interface {
	Exec() ([]LintOutput, error)
}
