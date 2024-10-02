// Package cmd provides the main commands for the gh-dxp extension.
package cmd

import (
	"context"

	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

// Execute runs the main command of the CLI tool.
func Execute(settings *config.Settings, version string) error {
	mainCmd := GenerateCmd(settings, version)
	ctx := context.Background()

	err := mainCmd.ExecuteContext(ctx)

	return err
}

// GenerateCmd sets up the command structure for the CLI tool using Cobra.
func GenerateCmd(settings *config.Settings, version string) *cobra.Command {
	var (
		debug bool
	)

	var retCmd = &cobra.Command{
		Use:           "dxp",
		Short:         "Extended Git & GitHub CLI workflows for linting, testing, code review and merges.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version,
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			log.DecreasePadding()
			log.SetLevel(log.InfoLevel)

			if debug {
				log.Info("Debug logs enabled")
				log.SetLevel(log.DebugLevel)
			}
		},
	}

	retCmd.PersistentFlags().BoolVar(&debug, "debug", false, "verbose logging")

	exe := utils.LinuxExecutor()

	retCmd.AddCommand(
		AliasCmd(exe),
		BranchCmd(exe),
		LintCmd(exe, settings),
		OwnerCmd(exe, settings),
		PRCmd(exe, settings),
		TestCmd(exe),
		TemplateCmd(exe, settings),
		StatusCmd(exe),
		UpgradeCmd(exe),
	)

	return retCmd
}
