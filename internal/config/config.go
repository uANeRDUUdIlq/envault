package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const defaultConfigFile = ".envault.json"

// Config holds the envault project configuration.
type Config struct {
	BackendURL  string   `json:"backend_url"`
	Project     string   `json:"project"`
	Recipients  []string `json:"recipients"`
	IdentityFile string  `json:"identity_file"`
}

// Load reads the config file from the given directory (or CWD if empty).
func Load(dir string) (*Config, error) {
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	path := filepath.Join(dir, defaultConfigFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save writes the config to the given directory.
func Save(dir string, cfg *Config) error {
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(dir, defaultConfigFile)
	return os.WriteFile(path, data, 0644)
}

// Validate returns an error if the config is missing required fields.
func Validate(cfg *Config) error {
	if cfg.Project == "" {
		return errors.New("config: project name is required")
	}
	if cfg.BackendURL == "" {
		return errors.New("config: backend_url is required")
	}
	if len(cfg.Recipients) == 0 {
		return errors.New("config: at least one recipient public key is required")
	}
	return nil
}

// ErrNotFound is returned when no config file exists in the directory.
var ErrNotFound = errors.New("config: .envault.json not found")
