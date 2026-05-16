package env

import (
	"fmt"
	"os"
)

// InjectFromFile parses a .env file and injects its variables into the
// current process environment using the provided options.
// It returns the injection result and any error encountered.
func InjectFromFile(path string, opts InjectOptions) (*InjectResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("inject_loader: read %s: %w", path, err)
	}
	vars, err := Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("inject_loader: parse %s: %w", path, err)
	}
	inj := NewInjector(opts)
	res, err := inj.Inject(vars)
	if err != nil {
		return nil, fmt.Errorf("inject_loader: inject: %w", err)
	}
	return res, nil
}

// EjectFromFile parses a .env file and removes its variables from the
// current process environment. Returns the list of keys that were removed.
func EjectFromFile(path string, opts InjectOptions) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("inject_loader: read %s: %w", path, err)
	}
	vars, err := Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("inject_loader: parse %s: %w", path, err)
	}
	inj := NewInjector(opts)
	return inj.Eject(vars), nil
}
