package cmd

import (
	"path/filepath"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/lint"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

// LintCmd creates a new command to run the linters defined in the .devxp config.
func LintCmd(exe utils.Executor, settings *config.Settings) *cobra.Command {
	opts := &lint.Options{}

	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Run linters on modified files in the repository.",
		Args:  cobra.MaximumNArgs(0),
		Long: heredoc.Docf(`
			Run linters on files in the repository. We use MegaLinter (an open-source lint aggregator) to run linting
			on the repository. By default, only files that have been modified in relation to the main branch are
			included in the lint.

			Some linters (e.g., prettier) provide auto-fix capabilities. To resolve linting errors automatically, use
			the --fix flag.
		`, "`"),
		Example: heredoc.Doc(`
			# Lint modified files in the repository
			$ gh dxp lint

			# Lint all files in repository
			$ gh dxp lint --all

			# Lint modified files in repository and fix errors
			$ gh dxp lint --fix

			# Lint all files in the current directory
			$ gh dxp lint -d .
		`),
		RunE: func(_ *cobra.Command, _ []string) error {

			currentPath, _ := filepath.Abs("./")

			err := utils.SetWorkDirToGitHubRoot(exe)
			if err != nil {
				return err
			}

			// Convert relative path provided in Directory into a path relative to the git root
			if opts.Directory != "" {
				fullPath := filepath.Join(currentPath, opts.Directory)

				gitRoot, err := utils.GetGitRootDirectory(exe)
				if err != nil {
					return err
				}

				opts.Directory = strings.TrimPrefix(fullPath, gitRoot)
			}

			return lint.Run(exe, settings, opts)
		},
	}

	fl := cmd.Flags()
	fl.BoolVarP(
		&opts.LintAll,
		"all",
		"a",
		false,
		"Lint all files in the repository",
	)
	fl.BoolVarP(
		&opts.Fix,
		"fix",
		"f",
		false,
		"Automatically fix linting errors",
	)
	fl.StringVarP(
		&opts.Directory,
		"directory",
		"d",
		"",
		"Lint all files under the given directory",
	)
	fl.StringVarP(
		&opts.LinterImage,
		"image",
		"i",
		"",
		"Specify which MegaLinter image to use.",
	)
	fl.StringVarP(
		&opts.Proxy,
		"proxy",
		"p",
		"",
		"Specify an https proxy server to be used in the MegaLinter container",
	)

	return cmd
}
