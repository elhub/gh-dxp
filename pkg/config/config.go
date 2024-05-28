package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

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

func DefaultSettings() *Settings {
	return &Settings{
		ProjectType: "default",
	}
}

func MergeSettings(source *Settings, user *Settings) *Settings {
	if user.ProjectType != "" {
		source.ProjectType = user.ProjectType
	}

	return source
}
