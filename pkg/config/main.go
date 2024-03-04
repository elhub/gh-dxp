package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ReadConfig(filepath string) (*Settings, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var cfg Settings
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
