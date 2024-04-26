package lint

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/elhub/gh-dxp/pkg/config"
	"github.com/elhub/gh-dxp/pkg/utils"
)

func DefaultLinters() map[string]Linter {
	return map[string]Linter{
		"golint":   GoLint{},
		"yamllint": YamlLint{},
	}
}

func Run(exe utils.Executor, settings *config.Settings, linters map[string]Linter) error {
	// iterate over settings.Linters and run each one
	var outputs []LinterOutput
	for _, lintEntry := range settings.Lint.Linters {
		if linter, ok := linters[lintEntry.Name]; ok {
			output, err := linter.Run(exe)
			outputs = append(outputs, output...)

			if err != nil {
				fmt.Printf("%s returned %d errors\n", lintEntry.Name, len(outputs))
			}
		} else {
			fmt.Printf("Linter %s not found\n", lintEntry.Name)
		}
	}

	// print the outputs
	for _, output := range outputs {
		fmt.Printf("%s:%d:%d: %s: %s\n", output.Path, output.Line, output.Column, output.Description, output.Severity)
	}

	return nil
}

func GetFiles(exe utils.Executor, separator string, extensions ...string) (string, error) {
	rootDir, rootErr := getRootDir(exe)
	if rootErr != nil {
		return "", rootErr
	}
	var files []string
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if stringInSlice(filepath.Ext(path), extensions) {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return "", err
	}
	fileString := strings.Join(files, separator)

	return fileString, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getRootDir(exe utils.Executor) (string, error) {
	output, err := exe.Command("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(output, "\n"), nil
}
