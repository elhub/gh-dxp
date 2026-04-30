package pr

import "github.com/charmbracelet/bubbles/table"

type PrAuthor = prAuthor
type PrRepository = prRepository

var InitialModel = initialModel //nolint:gochecknoglobals // Expose for testing

func (ui PullRequestUI) Rows() []table.Row {
	return ui.table.Rows()
}
