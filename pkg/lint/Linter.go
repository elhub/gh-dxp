package lint

type Linter interface {
	Exec() ([]lintOutput, error)
}
