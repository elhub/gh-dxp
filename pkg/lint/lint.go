package lint

import (
	"context"
	"fmt"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/utils"
)

func Run(exe utils.Executor, settings *config.Settings) error {
	// Run mega-linter-runner with the cupcake flavor.
	ctx := context.Background()

	err := exe.CommandContext(ctx, "npx", "mega-linter-runner", "--flavor", "cupcake")
	if err != nil {
		fmt.Printf("The Lint Process returned an error: %s\n", err)
		return err
	}

	return nil
}
