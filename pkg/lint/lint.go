// Package lint provides functions to run linting on the codebase.
package lint

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/ghutil"
	"github.com/elhub/gh-dxp/pkg/logger"
	"gopkg.in/yaml.v3"
)

const megaLinterConfigPath = ".mega-linter.yml"

// repoSetsGoToolchain reports whether the repository's own .mega-linter.yml (resolved relative to
// basePath) already defines a GOTOOLCHAIN value. When it does, gh-dxp should not force
// GOTOOLCHAIN=auto, since that would silently override the repository's pinned value (MegaLinter
// always prefers OS-level env vars over file config for the same key). If the file can't be read
// or parsed, a warning is logged so the failure stays visible, and false is returned (falling back
// to today's default behavior of forcing GOTOOLCHAIN=auto). basePath is accepted explicitly
// (rather than always reading the process's current working directory) so callers - including
// tests - can point at a specific directory without mutating global process state via os.Chdir.
func repoSetsGoToolchain(basePath string) bool {
	path := filepath.Join(basePath, megaLinterConfigPath)
	if !ghutil.FileExists(path) {
		return false
	}

	data, err := os.ReadFile(path)
	if err != nil {
		logger.Warnf("could not read %s to check for a pinned GOTOOLCHAIN: %s", path, err)
		return false
	}

	var cfg map[string]interface{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		logger.Warnf("could not parse %s to check for a pinned GOTOOLCHAIN: %s", path, err)
		return false
	}

	_, ok := cfg["GOTOOLCHAIN"]
	return ok
}

// Run runs the linting process using megalinter (https://github.com/oxsecurity/megalinter).
// Megalinter is an open-source linter aggregator that runs multiple linters in parallel. It requires NodeJS (npx) to be installed.
func Run(exe ghutil.Executor, settings *config.Settings, opts *Options) error {
	if opts == nil {
		opts = &Options{}
	}

	baseDir := opts.BaseDir
	if baseDir == "" {
		baseDir = "."
	}
	configPath := filepath.Join(baseDir, megaLinterConfigPath)

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
		args = append(args, "--image", settings.MegalinterImageVersion)
	} else {
		args = append(args, "--image", opts.LinterImage)
	}

	args = append(args, "-e", "LINTER_RULES_PATH=/tmp") // Prevents mega-linter from spamming lint configuration files into the repository
	if !repoSetsGoToolchain(baseDir) {
		// Allow go to upgrade to latest version when running the linter, preventing linting
		// errors. Skipped if the repo's own .mega-linter.yml already pins a GOTOOLCHAIN value,
		// since this flag would otherwise silently override it.
		args = append(args, "-e", "GOTOOLCHAIN=auto")
	}

	// Keep execution in the current process working directory, but force MegaLinter
	// to load config from the same source Run() uses for decision logic.
	if ghutil.FileExists(configPath) {
		absConfigPath, err := filepath.Abs(configPath)
		if err != nil {
			logger.Warnf("could not resolve absolute path for %s: %s", configPath, err)
			absConfigPath = configPath
		}
		args = append(args, "-e", "MEGALINTER_CONFIG="+absConfigPath)
	} else {
		logger.Info("Using the default Elhub mega-linter configuration.\n")
		// Append the default configuration file to the args.
		args = append(args, "-e", "MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml")
	}
	if !opts.LintAll && opts.Directory == "" {
		changedFiles, err := ghutil.GetChangedFiles(exe)
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
