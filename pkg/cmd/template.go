package cmd

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/template"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

// TemplateCmd initializes a repository with default files.
func TemplateCmd(_ utils.Executor, settings *config.Settings) *cobra.Command {
	opts := &template.Options{}

	cmd := &cobra.Command{
		Use:   "template",
		Short: "Copy standard template files to the current repository.",
		Long: heredoc.Docf(`
			Initialize a repository with default files defined in a project template. If files already exist, the
			standard files (i.e., those which should not be edited) will be overwritten.
		`, "`"),
		Example: heredoc.Doc(`
			$ gh dxp template
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, _ []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("could not get current working directory: %w", err)
			}

			return template.Execute(cwd, settings, opts)
		},
	}

	fl := cmd.Flags()
	fl.BoolVarP(
		&opts.IsGradleProject,
		"gradle",
		"m",
		false,
		"Include gradle specific files",
	)
	return cmd
}
