package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

func ReadConfig() (*Settings, error) {
	data, err := ioutil.ReadFile(".devxp")
	if err != nil {
		return nil, err
	}

	var cfg Settings
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
