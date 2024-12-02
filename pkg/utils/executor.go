// Package utils provides common utilities for the gh-dxp extension.
package utils

import (
	"context"
)

// Executor is an interface for running commands.
type Executor interface {
	Command(name string, args ...string) (string, error)
	CommandContext(ctx context.Context, name string, arg ...string) error
	GH(args ...string) (string, error)
	Chdir(dir string) error
}
