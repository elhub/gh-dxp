package pr

import "charm.land/bubbles/v2/table"

type PrAuthor = prAuthor
type PrRepository = prRepository

var InitialModel = initialModel           //nolint:gochecknoglobals // Expose for testing
var ValidateLocalChanges = validateLocalChanges //nolint:gochecknoglobals // Expose for testing
var AddAndCommitFiles = addAndCommitFiles

func (ui PullRequestUI) Rows() []table.Row {
	return ui.table.Rows()
}
