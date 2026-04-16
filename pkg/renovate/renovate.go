package renovate

import (
	"context"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/logger"
	"github.com/elhub/gh-dxp/pkg/utils"
)

func Run(exe utils.Executor, _ *config.Settings, opts *Options) error {
	err, renovateConfigChanged := isRenovateConfigUpdated(exe)
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

func isRenovateConfigUpdated(exe utils.Executor) (error, bool) {
	changedFiles, err := utils.GetChangedFiles(exe)
	if err != nil {
		return err, false
	}

	if len(changedFiles) == 0 {
		logger.Info("Did not find any changed files to validate")
		return nil, false
	}

	for _, file := range changedFiles {
		if file == ".github/renovate.json" {
			return nil, true
		}
	}
	return nil, false
}
