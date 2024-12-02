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

func initialModel(pullRequests []PullRequestInfo) PullRequestUI {

	rows := []table.Row{}
	//rows :=//make([]table.Row, 0, len(mine)+len(review)+1)
	var maxRepoNameLen = 12
	var maxTitleLen = 40
	var maxChangesLen = 8
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
		if len(pr.HeadRepository.Name) > maxRepoNameLen {
			maxRepoNameLen = len(pr.HeadRepository.Name)
		}
		if len(pr.Title) > maxTitleLen {
			maxTitleLen = len(pr.Title)
		}
		if len(changes) > maxChangesLen {
			maxChangesLen = len(changes)
		}
	}

	t := table.New(
		table.WithColumns([]table.Column{
			{Title: "Repository", Width: maxRepoNameLen},
			{Title: "ID", Width: 4},
			{Title: "Title", Width: maxTitleLen},
			{Title: "Author", Width: 30},
			{Title: "Status", Width: 16},
			{Title: "Changes", Width: maxChangesLen},
		}),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(pullRequests)),
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

func (ui PullRequestUI) Init() tea.Cmd {
	return nil
}

func (ui PullRequestUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return ui, tea.Quit
}

func (ui PullRequestUI) View() string {
	return baseStyle.Render(ui.table.View()) + "\n  " //+ ui.table.HelpView() + "\n"
}
