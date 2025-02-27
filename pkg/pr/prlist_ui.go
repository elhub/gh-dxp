package pr

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func initialModel(pullRequests []pullRequestInfo) pullRequestUI {

	rows := []table.Row{}
	var repoNameLen = 12
	var titleLen = 40
	var changesLen = 8
	for _, pr := range pullRequests {
		changes := "+" + strconv.Itoa(pr.Additions) + " -" + strconv.Itoa(pr.Deletions)
		rows = append(rows, table.Row{
			pr.HeadRepository.Name,
			"#" + strconv.Itoa(pr.Number),
			pr.Title,
			pr.Author.Name,
			pr.ReviewDecision,
			changes,
		})
	}

	t := table.New(
		table.WithColumns([]table.Column{
			{Title: "Repository", Width: repoNameLen},
			{Title: "ID", Width: 4},
			{Title: "Title", Width: titleLen},
			{Title: "Author", Width: 20},
			{Title: "Status", Width: 16},
			{Title: "Changes", Width: changesLen},
		}),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(pullRequests)+1),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Cell = s.Cell.
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(false)

	t.SetStyles(s)
	return pullRequestUI{
		table: t,
	}
}

func (ui pullRequestUI) Init() tea.Cmd {
	return nil
}

func (ui pullRequestUI) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return ui, tea.Quit
}

func (ui pullRequestUI) View() string {
	return baseStyle.Render(ui.table.View()) + "\n  "
}
