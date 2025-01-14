// Package template provides utilities to set up new repositories using our project template.
package template

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/elhub/gh-dxp/pkg/config"
)

// Execute downloads the project template files and writes them to the working directory.
func Execute(workingDir string, settings *config.Settings, options *Options) error {
	// Get the project template URI
	uri := settings.ProjectTemplateURI

	// Create the .github directory if it does not exist
	ghDir, err := createDirectory(workingDir, ".github")
	if err != nil {
		return err
	}

	// Create the .teamcity directory if it does not exist
	tcDir, err := createDirectory(workingDir, ".teamcity")
	if err != nil {
		return err
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
			path:      filepath.Join(ghDir, "CONTRIBUTING.md"),
			overwrite: true,
		},
		{
			fileName:  ".teamcity/pom-template.xml",
			path:      filepath.Join(tcDir, "pom.xml"),
			overwrite: false,
		},
		{
			fileName:  ".teamcity/settings-template.kts",
			path:      filepath.Join(tcDir, "settings.kts"),
			overwrite: false,
		},
	}

	if options.IsGradleProject {

		// Create gradle directories if we are setting up a gradle project
		gradleDir, err := createDirectory(workingDir, "gradle")
		if err != nil {
			return err
		}

		gradleWrapperDir, err := createDirectory(gradleDir, "wrapper")
		if err != nil {
			return err
		}

		gradleFiles := []struct {
			fileName  string
			path      string
			overwrite bool
		}{
			{
				fileName:  "gradleFiles/gradle/wrapper/gradle-wrapper.jar",
				path:      filepath.Join(gradleWrapperDir, "gradle-wrapper.jar"),
				overwrite: false,
			},
			{
				fileName:  "gradleFiles/gradle/wrapper/gradle-wrapper.properties",
				path:      filepath.Join(gradleWrapperDir, "gradle-wrapper.properties"),
				overwrite: false,
			},
			{
				fileName:  "gradleFiles/build.gradle.kts",
				path:      filepath.Join(workingDir, "build.gradle.kts"),
				overwrite: false,
			},
			{
				fileName:  "gradleFiles/gradle.properties",
				path:      filepath.Join(workingDir, "gradle.properties"),
				overwrite: false,
			},
			{
				fileName:  "gradleFiles/gradlew",
				path:      filepath.Join(workingDir, "gradlew"),
				overwrite: false,
			},
			{
				fileName:  "gradleFiles/settings.gradle.kts",
				path:      filepath.Join(workingDir, "settings.gradle.kts"),
				overwrite: false,
			},
		}

		files = append(files, gradleFiles...)
	}

	// Only write file if overwrite = true or file does not exist
	for _, file := range files {
		// Check if the file exists
		_, err := os.Stat(file.path)

		// If the file does not exist or overwrite is true, write the file
		if file.overwrite || os.IsNotExist(err) {
			err = writeFile(uri+file.fileName, file.path)
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
		}
	}

	return nil
}

// Downloads a file from an URI and writes it to path.
func writeFile(uri string, filepath string) error {
	// Create a new request
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, uri, nil)
	if err != nil {
		return err
	}

	// Create a new HTTP client and send the request
	resp, err := http.DefaultClient.Do(req)
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

func createDirectory(path string, name string) (string, error) {
	dir := filepath.Join(path, name)

	// Check if the directory already exists
	if _, errStat := os.Stat(dir); os.IsNotExist(errStat) {
		// Attempt to create the directory
		if err := os.Mkdir(dir, 0755); err != nil {
			return "", fmt.Errorf("could not create directory %s: %w", dir, err)
		}
	}
	return dir, nil
}
