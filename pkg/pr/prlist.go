// Package pr contains the functions and types for the pull request command.
package pr

import (
	"encoding/json"
	"sort"
	"strconv"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/elhub/gh-dxp/pkg/ghutil"
	"github.com/pkg/errors"
)

// ExecuteList renders the user's assigned pull requests.
func ExecuteList(exe ghutil.Executor, options *ListOptions) error {
	var wg sync.WaitGroup
	prChan := make(chan PullRequestInfo)
	errChan := make(chan error)

	if options.Mine {
		if err := retrievePullRequests("--author=@me", exe, prChan, errChan, &wg); err != nil {
			return err
		}
	}

	if options.ReviewRequested {
		if err := retrievePullRequests("--review-requested=@me", exe, prChan, errChan, &wg); err != nil {
			return err
		}
	}

	go func() {
		wg.Wait()
		close(prChan)
		close(errChan)
	}()

	pullRequests, err := drainChannels(prChan, errChan)
	if err != nil {
		return err
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

func retrievePullRequests(searchTerm string, exe ghutil.Executor, prChan chan<- PullRequestInfo, errChan chan<- error, wg *sync.WaitGroup) error {
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
		go fetchPullRequestDetails(exe, sr, prChan, errChan, wg)
	}

	return nil
}

func fetchPullRequestDetails(exe ghutil.Executor, sr searchResult, prChan chan<- PullRequestInfo, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	url := "https://github.com/" + sr.Repository.NameWithOwner + "/pull/" + strconv.Itoa(sr.Number)
	pullRequestDetails, err := exe.GH("pr", "view", url, "--json", "additions,author,createdAt,deletions,headRepository,number,title,reviewDecision")

	if err != nil {
		errChan <- errors.Wrap(err, "failed to get pr details")
		return
	}

	var pullRequest PullRequestInfo
	err = json.Unmarshal([]byte(pullRequestDetails), &pullRequest)
	if err != nil {
		errChan <- errors.Wrap(err, "failed to unmarshal pull request details")
		return
	}

	prChan <- pullRequest
}

func drainChannels(prChan <-chan PullRequestInfo, errChan <-chan error) ([]PullRequestInfo, error) {
	var result []PullRequestInfo
	var openPR, openErr = true, true
	for openPR || openErr {
		select {
		case pr, ok := <-prChan:
			if !ok {
				openPR = false
			} else {
				result = append(result, pr)
			}
		case err, ok := <-errChan:
			if ok {
				return nil, err
			}
			openErr = false
		}
	}
	return result, nil
}

func sortPullRequests(pullRequests []PullRequestInfo) {
	sort.Slice(pullRequests, func(i, j int) bool {
		if pullRequests[i].HeadRepository.ID == pullRequests[j].HeadRepository.ID {
			return pullRequests[i].Number < pullRequests[j].Number
		}
		return pullRequests[i].HeadRepository.Name < pullRequests[j].HeadRepository.Name
	})
}
