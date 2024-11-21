package pr

import (
	"encoding/json"
	"strconv"
	"sync"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cli/go-gh/v2"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/pkg/errors"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var pullRequests []PullRequestInfo
var wg sync.WaitGroup
var prChan chan PullRequestInfo
var errChan chan error

func ExecuteList(exe utils.Executor, options *ListOptions) error {
	pullRequests = []PullRequestInfo{}
	prChan = make(chan PullRequestInfo)
	errChan = make(chan error)

	if options.Mine {
		err := RetrievePullRequests("--author=@me")
		if err != nil {
			return err
		}
	}

	if options.ReviewRequested {
		RetrievePullRequests("--review-requested=@me")
	}

	go func() {
		wg.Wait()
		close(prChan)
		close(errChan)
	}()

	for prChan != nil || errChan != nil {
		select {
		case pr, ok := <-prChan:
			if !ok {
				prChan = nil
			} else {
				pullRequests = append(pullRequests, pr)
			}
		case err, ok := <-errChan:
			if ok {
				return err
			} else {
				errChan = nil
			}
		}
	}

	// Check content of pullRequests
	// Pretty Print myPullRequests as a Table
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		return err
	}

	/*
		// Process pullRequests as needed
		for _, pr := range pullRequests {
			log.Info(pr.Title)
		}*/

	return nil
}

func RetrievePullRequests(searchTerm string) error {

	res, _, err := gh.Exec("search", "prs", searchTerm, "--state=open", "--json", "number,repository")
	if err != nil {
		return errors.Wrap(err, "failed to search prs for my pull requests")
	}

	var searchResults []SearchResult
	err = json.Unmarshal(res.Bytes(), &searchResults)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal search results for my pull requests")
	}

	for _, sr := range searchResults {
		wg.Add(1)
		go fetchPullRequestDetails(sr, prChan, errChan, &wg)
	}

	return nil
}

func fetchPullRequestDetails(sr SearchResult, prChan chan<- PullRequestInfo, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	url := "https://github.com/" + sr.Repository.NameWithOwner + "/pull/" + strconv.Itoa(sr.Number)
	pullRequestDetails, _, err := gh.Exec("pr", "view", url, "--json", "additions,author,createdAt,deletions,headRepository,number,title,reviewDecision")
	if err != nil {
		errChan <- errors.Wrap(err, "failed to get pr details")
		return
	}

	var pullRequest PullRequestInfo
	err = json.Unmarshal(pullRequestDetails.Bytes(), &pullRequest)
	if err != nil {
		errChan <- errors.Wrap(err, "failed to unmarshal pull request details")
		return
	}

	prChan <- pullRequest
}

func initialModel() PullRequestUI {

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
