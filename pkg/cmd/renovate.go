// Package cmd provides the main commands for the gh-dxp extension.
package cmd

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/renovate"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

// RenovateCmd creates a new command for working with renovate.
func RenovateCmd(exe utils.Executor, settings *config.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "renovate",
		Short: "Work with renovate",
		Args:  cobra.MaximumNArgs(1),
		Long: heredoc.Doc(`
			The renovate command group allows you to validate renovate config.
		`),
	}

	cmd.AddCommand(ValidateCmd(exe, settings))

	return cmd
}

// ValidateCmd creates a new command for validating renovate config.
func ValidateCmd(exe utils.Executor, settings *config.Settings) *cobra.Command {
	opts := &renovate.Options{}
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "validate renovate config",
		Args:  cobra.MaximumNArgs(0),
		Long: heredoc.Docf(`
			Validate renovate config when changed.
		`, "`"),
		Example: heredoc.Doc(`
			# Validate renovate config if modified from 'main' branch
			$ gh dxp renovate validate

			# Force validation of renovate config (for example even if it's unchanged)
			$ gh dxp renovate validate --force
		`),
		RunE: func(_ *cobra.Command, _ []string) error {
			err := utils.SetWorkDirToGitHubRoot(exe)
			if err != nil {
				return err
			}
			return renovate.Run(exe, settings, opts)
		},
	}

	fl := cmd.Flags()
	fl.BoolVarP(
		&opts.Force,
		"force",
		"f",
		false,
		"Force validation even if there are no changes",
	)

	return cmd
}
