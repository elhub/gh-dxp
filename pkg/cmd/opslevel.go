// Package cmd provides CLI commands for the gh-dxp extension.
package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/elhub/gh-dxp/pkg/ghutil"
	"github.com/elhub/gh-dxp/pkg/opslevel"
	"github.com/spf13/cobra"
)

// OpsLevelCmd creates the command for generating and updating OpsLevel catalog configuration.
func OpsLevelCmd(exe ghutil.Executor) *cobra.Command {
	options := opslevel.Options{}
	cmd := &cobra.Command{
		Use:   "opslevel",
		Short: "Create an OpsLevel catalog configuration for the current repository.",
		Long: "Create opslevel.yml at the Git repository root using GitHub metadata.\n" +
			"Prompt for each metadata value with its current default, then review the YAML before writing.\n" +
			"Existing files require overwrite confirmation and provide default values.",
		Example: "  gh dxp opslevel\n" +
			"  gh dxp opslevel --no-input --type tool --lifecycle production --tier tier_4 --owner team-devxp --package-name no.elhub.my-tool",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, err := ghutil.GetGitRootDirectory(exe)
			if err != nil {
				return fmt.Errorf("could not locate repository root: %w", err)
			}
			if err := opslevel.Execute(root, exe, options); err != nil {
				if errors.Is(err, opslevel.ErrOverwriteDeclined) || errors.Is(err, opslevel.ErrWriteDeclined) {
					cmd.Println("No opslevel.yml changes written.")
					return nil
				}
				return err
			}
			cmd.Printf("Wrote %s\n", filepath.Join(root, "opslevel.yml"))
			return nil
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&options.Type, "type", "", "Component type: service, library, application, tool, configuration, orchestrator, test")
	flags.StringVar(&options.Lifecycle, "lifecycle", "", "Lifecycle: experimental, development, production, deprecated, retired")
	flags.StringVar(&options.Tier, "tier", "", "Criticality: tier_1, tier_2, tier_3, tier_4")
	flags.StringVar(&options.Owner, "owner", "", "OpsLevel team or chapter alias (not an individual)")
	flags.StringVar(&options.Description, "description", "", "Component description (defaults to existing value, then GitHub description)")
	flags.StringVar(&options.PackageName, "package-name", "", "Required package name alias for SonarQube integration")
	flags.StringVar(&options.System, "system", "", "OpsLevel system alias")
	flags.StringVar(&options.Language, "language", "", "Main language (defaults to existing value, then GitHub primary language)")
	flags.StringVar(&options.Framework, "framework", "", "Main framework or runtime")
	flags.BoolVar(&options.NoInput, "no-input", false, "Skip prompts and review; fail if required metadata is missing or the file exists")
	return cmd
}
