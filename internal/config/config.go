package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the gask configuration.
type Config struct {
	Editor string `json:"editor"`
}

var ErrConfigNotFound = errors.New("config file not found")

// Load reads the configuration from .gask/config.json in the current directory.
func Load() (*Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("config load: failed to get current directory: %w", err)
	}

	configPath := filepath.Join(cwd, ".gask", "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, ErrConfigNotFound
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("config load: failed to read file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config load: failed to unmarshal json: %w", err)
	}

	return &cfg, nil
}

// Save writes the configuration to .gask/config.json in the current directory.
func (c *Config) Save() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("config save: failed to get current directory: %w", err)
	}

	gaskDir := filepath.Join(cwd, ".gask")
	if _, err := os.Stat(gaskDir); os.IsNotExist(err) {
		if err := os.MkdirAll(gaskDir, 0755); err != nil {
			return fmt.Errorf("config save: failed to create .gask directory: %w", err)
		}
	}

	configPath := filepath.Join(gaskDir, "config.json")
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("config save: failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("config save: failed to encode json: %w", err)
	}

	return nil
}
