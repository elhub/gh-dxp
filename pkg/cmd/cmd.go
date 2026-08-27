// Package cmd provides the main commands for the gh-dxp extension.
package cmd

import (
	"context"
	"fmt"

	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/ghutil"
	"github.com/elhub/gh-dxp/pkg/jira"
	"github.com/elhub/gh-dxp/pkg/logger"
	"github.com/spf13/cobra"
)

// Execute runs the main command of the CLI tool.
func Execute(settings *config.Settings, version string) error {
	mainCmd, err := GenerateCmd(settings, version)
	if err != nil {
		return err
	}
	ctx := context.Background()

	err = mainCmd.ExecuteContext(ctx)

	return err
}

// GenerateCmd sets up the command structure for the CLI tool using Cobra.
func GenerateCmd(settings *config.Settings, version string) (*cobra.Command, error) {
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
			logger.DecreasePadding()
			logger.SetLevel(log.InfoLevel)

			if debug {
				logger.Info("Debug logs enabled")
				logger.SetLevel(log.DebugLevel)
			}
		},
	}

	retCmd.PersistentFlags().BoolVar(&debug, "debug", false, "verbose logging")

	exe := ghutil.LinuxExecutor()

	retCmd.AddCommand(
		AliasCmd(exe),
		BranchCmd(exe),
		JiraCmd(settings),
		LintCmd(exe, settings),
		OwnerCmd(exe, settings),
		PRCmd(exe, settings),
		RepoCmd(exe, settings),
		TestCmd(exe),
		TemplateCmd(exe, settings),
		StatusCmd(exe),
		RenovateCmd(exe, settings),
	)

	return retCmd, nil
}

func JiraCmd(settings *config.Settings) *cobra.Command {
	return &cobra.Command{
		Use:   "jira [commit message]",
		Short: "Suggest Jira issues for text",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			issues, err := jira.SearchIssues(cmd.Context(), settings.JiraURL, settings.JiraEmail, jira.SearchText{CommitMessage: args[0]})
			if err != nil {
				return err
			}
			for i, issue := range issues {
				if i == 5 {
					break
				}
				fmt.Printf("%d. %s - %s\n", i+1, issue.Key, issue.Fields.Summary)
			}
			return nil
		},
	}
}
