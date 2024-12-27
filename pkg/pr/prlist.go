package pr

import (
	"encoding/json"
	"sort"
	"strconv"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/pkg/errors"
)

var pullRequests []pullRequestInfo
var wg sync.WaitGroup
var prChan chan pullRequestInfo
var errChan chan error

// ExecuteList renders the user's assigned pull requests
func ExecuteList(exe utils.Executor, options *ListOptions) error {
	pullRequests = []pullRequestInfo{}
	prChan = make(chan pullRequestInfo)
	errChan = make(chan error)

	if options.Mine {
		err := retrievePullRequests("--author=@me", exe)
		if err != nil {
			return err
		}
	}

	if options.ReviewRequested {
		retrievePullRequests("--review-requested=@me", exe)
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
	sortPullRequests(pullRequests)

	// Check content of pullRequests
	// Pretty Print myPullRequests as a Table
	if !options.TestRun {
		p := tea.NewProgram(initialModel(pullRequests))
		if _, err := p.Run(); err != nil {
			return err
		}
	}

	return nil
}

func retrievePullRequests(searchTerm string, exe utils.Executor) error {

	res, err := exe.GH("search", "prs", searchTerm, "--state=open", "--json", "number,repository")
	if err != nil {
		return errors.Wrap(err, "failed to search prs for my pull requests")
	}

	var searchResults []searchResult
	err = json.Unmarshal([]byte(res), &searchResults)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal search results for my pull requests")
	}

	// Fetching the details of each PR is slow, so we do this in parallel
	for _, sr := range searchResults {
		wg.Add(1)
		go fetchPullRequestDetails(exe, sr, prChan, errChan, &wg)
	}

	return nil
}

func fetchPullRequestDetails(exe utils.Executor, sr searchResult, prChan chan<- pullRequestInfo, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	url := "https://github.com/" + sr.Repository.NameWithOwner + "/pull/" + strconv.Itoa(sr.Number)
	pullRequestDetails, err := exe.GH("pr", "view", url, "--json", "additions,author,createdAt,deletions,headRepository,number,title,reviewDecision")

	if err != nil {
		errChan <- errors.Wrap(err, "failed to get pr details")
		return
	}

	var pullRequest pullRequestInfo
	err = json.Unmarshal([]byte(pullRequestDetails), &pullRequest)
	if err != nil {
		errChan <- errors.Wrap(err, "failed to unmarshal pull request details")
		return
	}

	prChan <- pullRequest
}

func sortPullRequests(pullRequests []pullRequestInfo) {
	sort.Slice(pullRequests, func(i, j int) bool {
		if pullRequests[i].HeadRepository.ID == pullRequests[j].HeadRepository.ID {
			return pullRequests[i].Number < pullRequests[j].Number
		}
		return pullRequests[i].HeadRepository.Name < pullRequests[j].HeadRepository.Name
	})
}
