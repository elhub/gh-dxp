// Package init provides utilities to set up new repositories
package template

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/elhub/gh-dxp/pkg/config"
)

func Execute(workingDir string, settings *config.Settings) error {
	// Get the project template URI
	uri := settings.ProjectTemplateUri

	// Get the current working directory

	// Create the .github directory if it does not exist
	ghDir := filepath.Join(workingDir, ".github")
	if _, errStat := os.Stat(ghDir); os.IsNotExist(errStat) {
		if err := os.Mkdir(ghDir, 0755); err != nil {
			return fmt.Errorf("could not create .github directory: %w", err)
		}
	}

	// Download files
	files := []struct {
		fileName  string
		path      string
		overwrite bool
	}{
		{
			fileName:  ".editorconfig-template",
			path:      filepath.Join(workingDir, ".editorconfig"),
			overwrite: true,
		},
		{
			fileName:  ".gitattributes-template",
			path:      filepath.Join(workingDir, ".gitattributes"),
			overwrite: true,
		},
		{
			fileName:  ".gitignore-template",
			path:      filepath.Join(workingDir, ".gitignore"),
			overwrite: true,
		},
		{
			fileName:  "README-template.md",
			path:      filepath.Join(workingDir, "README.md"),
			overwrite: false,
		},
		{
			fileName:  ".github/CODEOWNERS-template",
			path:      filepath.Join(ghDir, "CODEOWNERS"),
			overwrite: false,
		},
		{
			fileName:  ".github/CONTRIBUTING-template.md",
			path:      filepath.Join(ghDir, "CONTRIBUTING"),
			overwrite: true,
		},
	}

	// Only write file if overwrite = true or file does not exist
	for _, file := range files {
		// Check if the file exists
		_, err := os.Stat(file.path)

		// If the file does not exist or overwrite is true, write the file
		if os.IsNotExist(err) || file.overwrite {
			err = writeFile(uri+file.fileName, file.path)
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
		}
	}

	return nil
}

// Downloads a file from an URI and writes it to path
func writeFile(uri string, filepath string) error {
	// Get the data
	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
