package pr

import (
	"encoding/json"
	"strconv"
	"sync"

	"github.com/caarlos0/log"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/pkg/errors"
)

var pullRequests []PullRequestInfo
var wg sync.WaitGroup
var prChan chan PullRequestInfo
var errChan chan error

func ExecuteList(exe utils.Executor, options *ListOptions) error {
	pullRequests = []PullRequestInfo{}
	prChan = make(chan PullRequestInfo)
	errChan = make(chan error)
	log.Error("Pre-wolo")
	if options.Mine {
		err := RetrievePullRequests("--author=@me", exe)
		if err != nil {
			return err
		}
	}

	if options.ReviewRequested {
		RetrievePullRequests("--review-requested=@me", exe)
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

	if !options.TestRun {
		p := tea.NewProgram(initialModel(pullRequests))
		if _, err := p.Run(); err != nil {
			return err
		}
	}

	/*
		// Process pullRequests as needed
		for _, pr := range pullRequests {
			log.Info(pr.Title)
		}*/

	return nil
}

func RetrievePullRequests(searchTerm string, exe utils.Executor) error {

	res, err := exe.GH("search", "prs", searchTerm, "--state=open", "--json", "number,repository")
	if err != nil {
		return errors.Wrap(err, "failed to search prs for my pull requests")
	}

	var searchResults []SearchResult
	err = json.Unmarshal(res.Bytes(), &searchResults)
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

func fetchPullRequestDetails(exe utils.Executor, sr SearchResult, prChan chan<- PullRequestInfo, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	url := "https://github.com/" + sr.Repository.NameWithOwner + "/pull/" + strconv.Itoa(sr.Number)
	pullRequestDetails, err := exe.GH("pr", "view", url, "--json", "additions,author,createdAt,deletions,headRepository,number,title,reviewDecision")

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
