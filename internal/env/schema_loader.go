package env

import (
	"fmt"
	"os"
	"path/filepath"
)

const defaultSchemaFile = ".envault-schema"

// LoadSchema reads a schema from the given path.
// If path is empty it looks for defaultSchemaFile in the working directory.
func LoadSchema(path string) (*Schema, error) {
	if path == "" {
		path = defaultSchemaFile
	}
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("open schema %q: %w", path, err)
	}
	defer f.Close()
	return ParseSchema(f)
}

// ValidateFile parses an env file and validates it against a schema file.
// It returns all violations found.
func ValidateFile(envPath, schemaPath string) ([]SchemaViolation, error) {
	ef, err := os.Open(filepath.Clean(envPath))
	if err != nil {
		return nil, fmt.Errorf("open env file: %w", err)
	}
	defer ef.Close()

	vars, err := Parse(ef)
	if err != nil {
		return nil, fmt.Errorf("parse env file: %w", err)
	}

	s, err := LoadSchema(schemaPath)
	if err != nil {
		return nil, err
	}

	return s.Validate(vars), nil
}
