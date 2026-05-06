package env

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ImportFormat represents a supported import file format.
type ImportFormat string

const (
	ImportFormatDotenv ImportFormat = "dotenv"
	ImportFormatJSON   ImportFormat = "json"
	ImportFormatShell  ImportFormat = "shell"
)

// ImportResult holds the result of an import operation.
type ImportResult struct {
	Vars     map[string]string
	Skipped  []string
	Imported int
}

// Importer reads env vars from external files into envault.
type Importer struct {
	format    ImportFormat
	overwrite bool
}

// NewImporter creates an Importer for the given format.
func NewImporter(format ImportFormat, overwrite bool) *Importer {
	return &Importer{format: format, overwrite: overwrite}
}

// ImportFile reads the file at path and returns an ImportResult.
func (i *Importer) ImportFile(path string, existing map[string]string) (*ImportResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("import: open file: %w", err)
	}
	defer f.Close()

	var parsed map[string]string
	switch i.format {
	case ImportFormatDotenv:
		parsed, err = parseDotenvReader(bufio.NewReader(f))
	case ImportFormatJSON:
		parsed, err = parseJSONReader(f)
	case ImportFormatShell:
		parsed, err = parseShellReader(bufio.NewReader(f))
	default:
		return nil, fmt.Errorf("import: unsupported format %q", i.format)
	}
	if err != nil {
		return nil, fmt.Errorf("import: parse: %w", err)
	}

	result := &ImportResult{Vars: make(map[string]string)}
	for k, v := range parsed {
		if _, exists := existing[k]; exists && !i.overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		result.Vars[k] = v
		result.Imported++
	}
	return result, nil
}

func parseDotenvReader(r *bufio.Reader) (map[string]string, error) {
	scanner := bufio.NewScanner(r)
	vars := make(map[string]string)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		k := strings.TrimSpace(parts[0])
		v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		if k != "" {
			vars[k] = v
		}
	}
	return vars, scanner.Err()
}

func parseJSONReader(f *os.File) (map[string]string, error) {
	vars := make(map[string]string)
	if err := json.NewDecoder(f).Decode(&vars); err != nil {
		return nil, err
	}
	return vars, nil
}

func parseShellReader(r *bufio.Reader) (map[string]string, error) {
	scanner := bufio.NewScanner(r)
	vars := make(map[string]string)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.TrimPrefix(line, "export ")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		k := strings.TrimSpace(parts[0])
		v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		if k != "" {
			vars[k] = v
		}
	}
	return vars, scanner.Err()
}
