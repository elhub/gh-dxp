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

	// Create the .github directory
	ghDir := filepath.Join(workingDir, ".github")
	if err := os.Mkdir(ghDir, 0755); err != nil {
		return fmt.Errorf("could not create .github directory: %w", err)
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
			overwrite: true,
		},
		{
			fileName:  "./github/CODEOWNERS-template",
			path:      filepath.Join(ghDir, "CODEOWNERS"),
			overwrite: true,
		},
		{
			fileName:  "./github/CONTRIBUTING-template.md",
			path:      filepath.Join(ghDir, "CONTRIBUTING"),
			overwrite: true,
		},
	}

	for _, f := range files {
		writeFile(uri+f.fileName, f.path)
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
