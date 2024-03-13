package lint

import "github.com/elhub/gh-devxp/pkg/utils"

type Linter interface {
	Exec(exec *utils.Executor) ([]LinterOutput, error)
}
