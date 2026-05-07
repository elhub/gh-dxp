// Package pr contains the functions and types for the pull request command.
package pr

import (
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"strconv"
)

func baseStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
}

func initialModel(pullRequests []PullRequestInfo) PullRequestUI {
	rows := []table.Row{}
	var repoNameLen = 12
	var titleLen = 40
	var changesLen = 8
	cols := []table.Column{
		{Title: "Repository", Width: repoNameLen},
		{Title: "ID", Width: 4},
		{Title: "Title", Width: titleLen},
		{Title: "Author", Width: 20},
		{Title: "Status", Width: 16},
		{Title: "Changes", Width: changesLen},
	}
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

	totalWidth := 0
	for _, col := range cols {
		totalWidth += col.Width + 2
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(pullRequests)+1),
		table.WithWidth(totalWidth),
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
	return PullRequestUI{
		table: t,
	}
}

// Init is the initial command for the Bubble Tea program.
func (ui PullRequestUI) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model accordingly.
func (ui PullRequestUI) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return ui, tea.Quit
}

// View renders the UI.
func (ui PullRequestUI) View() tea.View {
	return tea.NewView(baseStyle().Render(ui.table.View()) + "\n  ")
}
