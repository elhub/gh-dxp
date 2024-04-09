package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	tests := []struct {
		name        string
		expectPanic bool
	}{
		{
			name:        "Test when cmd.Execute succeeds",
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		if tt.expectPanic {
			assert.Panics(t, main, tt.name)
		} else {
			assert.NotPanics(t, main, tt.name)
		}
	}
}
