package pr

import (
	"github.com/caarlos0/log"
	"github.com/cli/go-gh/v2"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/pkg/errors"
)

func ExecuteList(exe utils.Executor, options *ListOptions) error {

	// Show the list of my PRs
	// gh search prs --author=@me --state=open --json number,title,author,headRefName,baseRefName
	if options.Mine {
		log.Info("Listing your PRs")
		res, _, err := gh.Exec("search", "prs", "--author=@me", "--state=open")

		if err != nil {
			return errors.Wrap(err, "failed to list PRs")
		}
		log.Info(res.String())
	}

	return nil
}
