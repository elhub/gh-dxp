// Package lint provides functions to run linting on the codebase.
package lint

import (
	"context"
	"fmt"
	"strings"

	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/utils"
)

// Run runs the linting process using megalinter (https://github.com/oxsecurity/megalinter).
// Megalinter is an open-source linter aggregator that runs multiple linters in parallel. It requires NodeJS (npx) to be installed.
func Run(exe utils.Executor, _ *config.Settings, opts *Options) error {
	// Run mega-linter-runner with the cupcake flavor.
	ctx := context.Background()

	args := []string{"npx", "mega-linter-runner"}

	if opts.LinterImage == "" {
		args = append(args, "--flavor", "cupcake")
	} else {
		args = append(args, "--image", opts.LinterImage)
	}

	// Check if mega-lint configuration is present in the repository.
	if !utils.FileExists(".mega-linter.yml") {
		log.Info("Using the default Elhub mega-linter configuration.\n")
		// Append the default configuration file to the args.
		args = append(args, "-e", "MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml")
	}
	if !opts.LintAll && opts.Directory == "" {
		changedFiles, err := getChangedFiles(exe)
		if err != nil {
			return err
		}

		if len(changedFiles) == 0 {
			log.Info("Did not find any changed files to lint")
			return nil
		}

		args = append(args, "--filesonly")
		args = append(args, changedFiles...)
	} else if opts.Directory != "" {
		args = append(args, "-e", "FILTER_REGEX_INCLUDE="+fmt.Sprintf("(%s)", opts.Directory))
	}
	if opts.Fix {
		args = append(args, "--fix")
	}
	err := exe.CommandContext(ctx, args[0], args[1:]...)
	if err != nil {
		log.Info("The Lint Process returned an error: " + err.Error() + "\n")
		return err
	}
	return nil
}

func getChangedFiles(exe utils.Executor) ([]string, error) {
	branchString, err := exe.Command("git", "branch")
	if err != nil {
		return []string{}, err
	}

	branchList := ConvertTerminalOutputIntoList(branchString)

	var changedFiles []string

	if len(branchList) > 0 {
		changedFilesString, err := exe.Command("git", "diff", "--name-only", "main", "--relative")
		changedFiles = ConvertTerminalOutputIntoList(changedFilesString)
		if err != nil {
			return []string{}, err
		}
	} else {
		changedFiles, err = utils.GetTrackedChanges(exe)
		if err != nil {
			return []string{}, err
		}
	}
	return changedFiles, nil
}

// ConvertTerminalOutputIntoList converts terminal output on multiple lines to a list of strings
func ConvertTerminalOutputIntoList(changedFilesString string) []string {
	if len(changedFilesString) == 0 {
		return []string{}
	}
	return strings.Split(strings.TrimSpace(changedFilesString), "\n")
}
