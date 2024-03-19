package config_test

import (
	"os"
	"testing"

	"github.com/elhub/gh-devxp/pkg/config"
)

func TestReadConfig(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", ".devxp")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write some data to the file
	text := []byte(`---
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
    - "(\\.break$)"`)
	if _, err1 := tmpfile.Write(text); err1 != nil {
		t.Fatal(err1)
	}
	if err2 := tmpfile.Close(); err2 != nil {
		t.Fatal(err2)
	}

	// Set the file name in ReadConfig
	cfg, err := config.ReadConfig(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Check that the settings were correctly read
	if len(cfg.Lint.Linters) != 2 {
		t.Errorf("got %+v, want 2", cfg.Lint.Linters)
	}

	if cfg.Lint.Linters[0].Name != "linter1" {
		t.Errorf("got %+v, want linter1", cfg)
	}

	if cfg.Lint.Linters[1].Name != "linter2" {
		t.Errorf("got %+v, want linter2", cfg)
	}

	if cfg.Lint.Linters[1].Include[0] != ".*\\.txt$" {
		t.Errorf("got %+v, want linter2 include", cfg)
	}

	if cfg.Lint.Linters[0].Exclude[0] != ".*\\.go$" {
		t.Errorf("got %+v, want linter1 exclude", cfg)
	}

	if len(cfg.Lint.Exclude) != 2 {
		t.Errorf("got %+v, want 2", cfg.Lint.Exclude)
	}
}

func TestFailConfig(t *testing.T) {
	// Check that an incorrect file returns a error
	_, err := config.ReadConfig(".devxpp")
	if err == nil {
		t.Fatal()
	}
}

func TestBadYamlConfig(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", ".failxp")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write some incorrectly formatted YAML to the file
	text := []byte(`---
lint:
  linters:
    - navn: linter1
	 - name: linter2`)
	if _, err1 := tmpfile.Write(text); err1 != nil {
		t.Fatal(err1)
	}
	if err2 := tmpfile.Close(); err2 != nil {
		t.Fatal(err2)
	}

	_, errFail := config.ReadConfig(tmpfile.Name())
	if errFail == nil {
		t.Fatal(err)
	}
}
