package cmd

import (
	"context"
	"os"

	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

func Execute(version string) {
	mainCmd := GenerateCmd(version)
	ctx := context.Background()

	if err := mainCmd.ExecuteContext(ctx); err != nil {
		log.WithError(err).Error("Command failed")
		os.Exit(1)
	}
}

func GenerateCmd(version string) *cobra.Command {
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

	retCmd.AddCommand(
		WorkCmd(),
	)

	return retCmd
}
