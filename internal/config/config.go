package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

func LoadConfig(path string) (*Config, error) {

	var err error
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	// Check if file is empty or contains only whitespace
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, fmt.Errorf("config file is empty")
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
