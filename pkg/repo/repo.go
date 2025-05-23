package repo

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/elhub/gh-dxp/pkg/logger"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/pkg/errors"
)

type repositoryInfo struct {
	Name     string `json:"name"`
	FullName string `json:"fullName"`
	URL      string `json:"url"`
}

// ExecuteClone carries out a clone all on the given pattern
func ExecuteClone(exe utils.Executor, pattern string, sleepFunction func(time.Duration), opts *Options) error {
	repositoryInfo, err := retrieveRepositories(pattern, exe)
	if err != nil {
		return err
	}

	for _, repo := range repositoryInfo {
		// If directory exits, skip cloning
		exists, err := utils.DirectoryExists(repo.Name)
		if err != nil {
			return errors.Wrap(err, "failed to check if directory exists")
		}

		if !exists {
			if opts.DryRun {
				logger.Infof("Dry run: Clone repository %s", repo.FullName)
			} else {
				logger.Infof("Cloning repository: %s", repo.FullName)

				err := cloneRepoWithRetries(repo.FullName, sleepFunction, exe)
				if err != nil {
					logger.Warn(err.Error())
				}
			}
		} else {
			logger.Infof("Skipped cloning of repository %s as directory already exists", repo.FullName)
		}
	}

	return nil
}

func retrieveUserOrgs(exe utils.Executor) ([]string, error) {
	type Org struct {
		Login string `json:"login"`
	}

	res, err := exe.GH("api", "user/orgs")
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve user organizations")
	}

	var orgs []Org
	err = json.Unmarshal([]byte(res), &orgs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal user organizations")
	}

	var orgLogins []string
	for _, org := range orgs {
		orgLogins = append(orgLogins, org.Login)
	}
	return orgLogins, nil
}

func retrieveRepositories(pattern string, exe utils.Executor) ([]repositoryInfo, error) {
	orgs, err := retrieveUserOrgs(exe)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve user organizations")
	}
	// Convert to comma separated string
	orgsStr := strings.Join(orgs, ",")
	logger.Infof("Searching for repositories in the following organizations: %s", orgsStr)

	var res string
	if len(pattern) == 0 {
		res, err = exe.GH("search", "repos", "--archived=false", "--json", "name,fullName,url", "--limit=1000", "--owner", orgsStr)
	} else {
		res, err = exe.GH("search", "repos", pattern, "--match", "name", "--archived=false", "--json", "name,fullName,url", "--limit=1000", "--owner", orgsStr)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to search repositories")
	}

	var searchResults []repositoryInfo
	err = json.Unmarshal([]byte(res), &searchResults)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal search results for repositories")
	}
	logger.Infof("Found %d repositories matching the pattern %s", len(searchResults), pattern)

	return searchResults, nil
}

func cloneRepoWithRetries(reponame string, sleep func(time.Duration), exe utils.Executor) error {

	maxAttempts := 5

	for i := 0; i <= maxAttempts; i++ {
		_, err := exe.GH("repo", "clone", reponame)
		if err != nil {
			if i != maxAttempts {
				sleepDuration := powInt(2, i)
				logger.Debugf("Will retry clone in %d seconds", sleepDuration)
				sleep(time.Second * time.Duration(sleepDuration))
			}
			continue
		}
		return nil
	}
	return fmt.Errorf("unable to clone repo %s after %d attempts", reponame, maxAttempts)
}

func powInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}
