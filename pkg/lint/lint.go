// Package lint provides functions to run linting on the codebase.
package lint

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/logger"
	"github.com/elhub/gh-dxp/pkg/utils"
)

// Run runs the linting process using megalinter (https://github.com/oxsecurity/megalinter).
// Megalinter is an open-source linter aggregator that runs multiple linters in parallel. It requires NodeJS (npx) to be installed.
func Run(exe utils.Executor, _ *config.Settings, opts *Options) error {
	// Create a context that listens for interrupt signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)
	go func() {
		select {
		case <-sigChan:
			logger.Info("\nReceived interrupt signal, stopping linter...\n")
			cancel()
		case <-ctx.Done():
			// Context cancelled, exit goroutine
			return
		}
	}()

	args := []string{"npx", "mega-linter-runner"}

	if opts.LinterImage == "" {
		args = append(args, "--flavor", "cupcake")
	} else {
		args = append(args, "--image", opts.LinterImage)
	}

	args = append(args, "-e", "LINTER_RULES_PATH=/tmp") // Prevents mega-linter from spamming lint configuration files into the repository

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
