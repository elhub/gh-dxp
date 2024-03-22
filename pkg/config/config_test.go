package config_test

import (
	"os"
	"testing"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTempFile(t *testing.T, text []byte) *os.File {
	tmpfile, err := os.CreateTemp("", ".devxp")
	if err != nil {
		t.Fatal(err)
	}

	if _, err1 := tmpfile.Write(text); err1 != nil {
		t.Fatal(err1)
	}
	if err2 := tmpfile.Close(); err2 != nil {
		t.Fatal(err2)
	}

	return tmpfile
}

func TestReadConfig(t *testing.T) {
	t.Run("valid config file", func(t *testing.T) {
		// Create a temporary file
		tmpfile := writeTempFile(t, []byte(`---
lint:
  linters:
    - name: linter1
      exclude:
        - ".*\\.go$"
    - name: linter2
      include:
        - ".*\\.txt$"
  exclude:
    - "(\\.bad$)"
    - "(\\.break$)"`))

		// Read the tmpfile
		cfg, err := config.ReadConfig(tmpfile.Name())

		// Check that the settings were correctly read
		require.NoError(t, err)
		assert.Len(t, cfg.Lint.Linters, 2)
		assert.Equal(t, "linter1", cfg.Lint.Linters[0].Name)
		assert.Equal(t, "linter2", cfg.Lint.Linters[1].Name)
		assert.Equal(t, ".*\\.txt$", cfg.Lint.Linters[1].Include[0])
		assert.Equal(t, ".*\\.go$", cfg.Lint.Linters[0].Exclude[0])
		assert.Len(t, cfg.Lint.Exclude, 2)
	})

	t.Run("non existent config file", func(t *testing.T) {
		_, err := config.ReadConfig(".devxpp")
		require.Error(t, err)
	})

	t.Run("incorrectly formatted YAML", func(t *testing.T) {
		// Create a temporary file
		tmpfile := writeTempFile(t, []byte(`---
lint:
  linters:
    - navn: linter1
	- name: linter2`))

		// Read the tmpfile
		_, err := config.ReadConfig(tmpfile.Name())

		// Check that the settings read throws an error
		require.Error(t, err)
	})
}
