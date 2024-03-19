package lint

import "github.com/elhub/gh-dxp/pkg/utils"

type Linter interface {
	Exec(exec *utils.Executor) ([]LinterOutput, error)
}
