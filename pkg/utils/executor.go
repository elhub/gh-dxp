// Package utils provides common utilities for the gh-dxp extension.
package utils

import (
	"bytes"
	"context"
)

// Executor is an interface for running commands.
type Executor interface {
	Command(name string, args ...string) (string, error)
	CommandContext(ctx context.Context, name string, arg ...string) error
	GH(args ...string) (bytes.Buffer, error)
	Chdir(dir string) error
}
