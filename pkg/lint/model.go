package lint

import "github.com/elhub/gh-dxp/pkg/utils"

type Linter interface {
	Exec(exec *utils.Executor) ([]LinterOutput, error)
}

type LinterOutput struct {
	Linter      string
	Path        string
	Line        int
	Column      int
	Description string
	Severity    string
	Source      string
}
