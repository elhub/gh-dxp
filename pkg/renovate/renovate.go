// Package renovate provides functionality to validate renovate configuration files.
package renovate

import (
	"context"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/ghutil"
	"github.com/elhub/gh-dxp/pkg/logger"
)

// Run executes the renovate validation process.
func Run(exe ghutil.Executor, _ *config.Settings, opts *Options) error {
	renovateConfigChanged, err := isRenovateConfigUpdated(exe)
	if err != nil {
		logger.Info("The validation process returned an error looking for renovate config file: " + err.Error() + "\n")
		return err
	}

	if !renovateConfigChanged && !opts.Force {
		logger.Info("The renovate config has not changed, skipping...")
		return nil
	}

	args := []string{"npx", "--package", "renovate@43.78.0", "renovate-config-validator", "--strict"}
	ctx := context.Background()
	err = exe.CommandContext(ctx, args[0], args[1:]...)

	if err != nil {
		logger.Info("The validation process returned an error: " + err.Error() + "\n")
		return err
	}
	return nil
}

// isRenovateConfigUpdated checks if the renovate config file has been updated compared to the main branch.
func isRenovateConfigUpdated(exe ghutil.Executor) (bool, error) {
	changedFiles, err := ghutil.GetChangedFiles(exe)
	if err != nil {
		return false, err
	}

	if len(changedFiles) == 0 {
		logger.Info("Did not find any changed files to validate")
		return false, nil
	}

	for _, file := range changedFiles {
		if file == ".github/renovate.json" {
			return true, nil
		}
	}
	return false, nil
}
