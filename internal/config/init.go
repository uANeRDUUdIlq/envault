package config

import (
	"errors"
	"fmt"
	"os"
)

// InitOptions holds parameters for initialising a new envault project.
type InitOptions struct {
	Dir          string
	Project      string
	BackendURL   string
	Recipients   []string
	IdentityFile string
	Overwrite    bool
}

// Init creates a new .envault.json in the given directory.
// It returns ErrAlreadyExists if a config file already exists and Overwrite is false.
func Init(opts InitOptions) (*Config, error) {
	if opts.Dir == "" {
		var err error
		opts.Dir, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("config: resolving working directory: %w", err)
		}
	}

	_, err := Load(opts.Dir)
	if err == nil && !opts.Overwrite {
		return nil, ErrAlreadyExists
	}
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, fmt.Errorf("config: checking existing config: %w", err)
	}

	cfg := &Config{
		BackendURL:   opts.BackendURL,
		Project:      opts.Project,
		Recipients:   opts.Recipients,
		IdentityFile: opts.IdentityFile,
	}
	if err := Validate(cfg); err != nil {
		return nil, err
	}
	if err := Save(opts.Dir, cfg); err != nil {
		return nil, fmt.Errorf("config: saving new config: %w", err)
	}
	return cfg, nil
}

// ErrAlreadyExists is returned when Init is called in a directory that already has a config.
var ErrAlreadyExists = errors.New("config: .envault.json already exists; use --overwrite to replace it")
