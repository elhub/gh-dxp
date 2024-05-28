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

// TestReadConfig tests the ReadConfig function
func TestReadConfig(t *testing.T) {
	t.Run("valid config file", func(t *testing.T) {
		// Create a temporary file
		tmpfile := writeTempFile(t, []byte(`---
projectType: "go"`))

		// Read the tmpfile
		cfg, err := config.ReadConfig(tmpfile.Name())

		// Check that the settings were correctly read
		require.NoError(t, err)
		assert.Equal(t, "go", cfg.ProjectType)
	})

	t.Run("non existent config file", func(t *testing.T) {
		_, err := config.ReadConfig(".devxpp")
		require.Error(t, err)
	})

	t.Run("incorrectly formatted YAML", func(t *testing.T) {
		// Create a temporary file
		tmpfile := writeTempFile(t, []byte(`---
projectTypes:
	"go"`))

		// Read the tmpfile
		_, err := config.ReadConfig(tmpfile.Name())

		// Check that the settings read throws an error
		require.Error(t, err)
	})
}

// Test the MergeSettings function
func TestMergeSettings(t *testing.T) {
	defaultSettings := config.DefaultSettings()
	userSettings := &config.Settings{
		ProjectType: "go",
	}

	mergedSettings := config.MergeSettings(defaultSettings, userSettings)

	assert.Equal(t, "go", mergedSettings.ProjectType)
}
