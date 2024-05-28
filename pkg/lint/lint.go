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

	err := exe.CommandContext(ctx, "npx", "mega-linter-runner", "--flavor", "cupcake")
	if err != nil {
		log.Info("The Lint Process returned an error: " + err.Error() + "\n")
		return err
	}

	return nil
}
