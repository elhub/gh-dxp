// Package lint provides functions to run linting on the codebase.
package lint

import (
	"context"
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

	args := []string{"npx", "mega-linter-runner", "--flavor", "cupcake"}

	// Check if mega-lint configuration is present in the repository.
	if !utils.FileExists(".mega-linter.yml") {
		log.Info("Using the default Elhub mega-linter configuration.\n")
		// Append the default configuration file to the args.
		args = append(args, "-e", "MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml")
	}
	if !opts.LintAll {
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
	}

	err := exe.CommandContext(ctx, args[0], args[1:]...)
	if err != nil {
		log.Info("The Lint Process returned an error: " + err.Error() + "\n")
		return err
	}
	return nil
}

func getChangedFiles(exe utils.Executor) ([]string, error) {
	changedFilesString, err := exe.Command("git", "diff", "--name-only", "main", "--relative")
	if err != nil {
		return []string{}, err
	}

	return ConvertChangedFilesIntoList(changedFilesString), nil
}

// ConvertChangedFilesIntoList converts the output string of a git diff --name-only into a list of file paths.
func ConvertChangedFilesIntoList(changedFilesString string) []string {
	if len(changedFilesString) == 0 {
		return []string{}
	}
	return strings.Split(strings.TrimSpace(changedFilesString), "\n")
}
