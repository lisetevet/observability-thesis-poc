package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type MongoDB struct {
	URI        string `yaml:"uri"`
	Database   string `yaml:"database"`
	Collection string `yaml:"collection"`
}

type Config struct {
	Port     int    `yaml:"port"`
	BasePath string `yaml:"base_path"`
	MongoDB     MongoDB     `yaml:"mongodb"`
	UsersServiceURL string `yaml:"users_service_url"`
}

func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}