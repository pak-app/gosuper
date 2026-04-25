package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (*Config, error) {

	var err error
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
