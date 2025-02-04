// Package lint provides functions to run linting on the codebase.
package lint

import (
	"context"
	"fmt"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/logger"
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
		logger.Info("Using the default Elhub mega-linter configuration.\n")
		// Append the default configuration file to the args.
		args = append(args, "-e", "MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml")
	}
	if !opts.LintAll && opts.Directory == "" {
		changedFiles, err := utils.GetChangedFiles(exe)
		if err != nil {
			return err
		}

		if len(changedFiles) == 0 {
			logger.Info("Did not find any changed files to lint")
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
	if opts.Proxy != "" {
		args = append(args, "-e", fmt.Sprintf("https_proxy=%s", opts.Proxy))
	}
	err := exe.CommandContext(ctx, args[0], args[1:]...)
	if err != nil {
		logger.Info("The Lint Process returned an error: " + err.Error() + "\n")
		return err
	}
	return nil
}
