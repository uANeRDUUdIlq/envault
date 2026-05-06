package env

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ParseSchema reads a simple schema definition from r.
//
// Format (one field per line):
//
//	# comment
//	KEY [required] [pattern=<regex>] [# description]
func ParseSchema(r io.Reader) (*Schema, error) {
	schema := &Schema{}
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Strip inline comment
		description := ""
		if idx := strings.Index(line, " # "); idx != -1 {
			description = strings.TrimSpace(line[idx+3:])
			line = strings.TrimSpace(line[:idx])
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		field := SchemaField{
			Key:         parts[0],
			Description: description,
		}

		for _, tok := range parts[1:] {
			switch {
			case tok == "required":
				field.Required = true
			case strings.HasPrefix(tok, "pattern="):
				field.Pattern = strings.TrimPrefix(tok, "pattern=")
			default:
				return nil, fmt.Errorf("line %d: unknown token %q", lineNum, tok)
			}
		}

		schema.Fields = append(schema.Fields, field)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return schema, nil
}
