// Package lint provides functions to run linting on the codebase.
package lint

import (
	"context"

	"github.com/caarlos0/log"
	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/utils"
)

// Run runs the linting process using megalinter (https://github.com/oxsecurity/megalinter).
// Megalinter is an open-source linter aggregator that runs multiple linters in parallel. It requires NodeJS (npx) to be installed.
func Run(exe utils.Executor, _ *config.Settings) error {
	// Run mega-linter-runner with the cupcake flavor.
	ctx := context.Background()

	args := []string{"npx", "mega-linter-runner", "--flavor", "cupcake"}

	// Check if mega-lint configuration is present in the repository.
	if !utils.FileExists(".mega-linter.yml") {
		log.Info("Using the default Elhub mega-linter configuration.\n")
		// Append the default configuration file to the args.
		args = append(args, "-e", "MEGALINTER_CONFIG=https://raw.githubusercontent.com/elhub/devxp-lint-configuration/main/resources/.mega-linter.yml")
	}

	err := exe.CommandContext(ctx, args[0], args[1:]...)
	if err != nil {
		log.Info("The Lint Process returned an error: " + err.Error() + "\n")
		return err
	}
	// say some thing funny
	return nil
}
