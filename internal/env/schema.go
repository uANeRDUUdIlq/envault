package env

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaField describes a single expected env variable.
type SchemaField struct {
	Key         string
	Required    bool
	Description string
	Pattern     string // optional regex the value must match
}

// Schema holds a set of field definitions for an env file.
type Schema struct {
	Fields []SchemaField
}

// SchemaViolation represents a single validation failure.
type SchemaViolation struct {
	Key     string
	Message string
}

func (v SchemaViolation) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// Validate checks vars against the schema and returns all violations.
func (s *Schema) Validate(vars map[string]string) []SchemaViolation {
	var violations []SchemaViolation

	for _, field := range s.Fields {
		val, ok := vars[field.Key]
		if !ok || strings.TrimSpace(val) == "" {
			if field.Required {
				violations = append(violations, SchemaViolation{
					Key:     field.Key,
					Message: "required key is missing or empty",
				})
			}
			continue
		}
		if field.Pattern != "" {
			re, err := regexp.Compile(field.Pattern)
			if err != nil {
				violations = append(violations, SchemaViolation{
					Key:     field.Key,
					Message: fmt.Sprintf("invalid pattern %q: %v", field.Pattern, err),
				})
				continue
			}
			if !re.MatchString(val) {
				violations = append(violations, SchemaViolation{
					Key:     field.Key,
					Message: fmt.Sprintf("value %q does not match pattern %q", val, field.Pattern),
				})
			}
		}
	}
	return violations
}

// IsValid returns true when Validate produces no violations.
func (s *Schema) IsValid(vars map[string]string) bool {
	return len(s.Validate(vars)) == 0
}
