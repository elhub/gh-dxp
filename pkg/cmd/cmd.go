package cmd

import (
	"context"

	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/utils"
	"github.com/spf13/cobra"
)

func Execute(settings *config.Settings, version string) error {
	mainCmd := GenerateCmd(settings, version)
	ctx := context.Background()

	err := mainCmd.ExecuteContext(ctx)

	return err
}

func GenerateCmd(settings *config.Settings, version string) *cobra.Command {
	var (
		debug bool
	)

	var retCmd = &cobra.Command{
		Use:           "devxp",
		Short:         "Extended Git & GitHub CLI workflows",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
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
		PRCmd(exe),
		LintCmd(exe, settings),
		MergeCmd(exe),
		BranchCmd(exe),
	)

	return retCmd
}
