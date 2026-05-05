package env

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Entry represents a single key-value pair from a .env file.
type Entry struct {
	Key   string
	Value string
}

// Parse reads a .env file format from r and returns the parsed entries.
// Lines beginning with '#' are treated as comments and skipped.
// Empty lines are skipped. Lines must be in KEY=VALUE format.
func Parse(r io.Reader) ([]Entry, error) {
	var entries []Entry
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			return nil, fmt.Errorf("line %d: missing '=' in %q", lineNum, line)
		}

		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])

		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", lineNum)
		}

		// Strip surrounding quotes from value if present.
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		entries = append(entries, Entry{Key: key, Value: value})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning input: %w", err)
	}

	return entries, nil
}

// Serialize writes entries back into .env file format.
func Serialize(entries []Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "%s=%s\n", e.Key, e.Value)
	}
	return sb.String()
}
