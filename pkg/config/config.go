// Package config provides the logic for reading the configuration settings from a file.
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ReadConfig reads the configuration settings from the specified file and unmarshals it into a Settings struct.
func ReadConfig(filepath string) (*Settings, error) {
	data, readErr := os.ReadFile(filepath)
	if readErr != nil {
		return nil, readErr
	}

	var cfg Settings
	if yamlErr := yaml.Unmarshal(data, &cfg); yamlErr != nil {
		return nil, yamlErr
	}

	return &cfg, nil
}

// DefaultSettings loads the default .devxp settings.
func DefaultSettings() *Settings {
	return &Settings{
		ProjectTemplateUri: "https://raw.githubusercontent.com/elhub/devxp-project-template/main/resources/",
		ProjectType:        "",
	}
}

// MergeSettings merges two settings.
func MergeSettings(source *Settings, newSettings *Settings) *Settings {
	if newSettings.ProjectType != "" {
		source.ProjectTemplateUri = newSettings.ProjectTemplateUri
		source.ProjectType = newSettings.ProjectType
	}

	return source
}
