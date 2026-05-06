package env

import (
	"fmt"
	"strings"
)

// TemplateVar represents a variable declared in a template file.
type TemplateVar struct {
	Key      string
	Default  string
	Required bool
	Comment  string
}

// Template holds the parsed structure of a .env.template file.
type Template struct {
	Vars []TemplateVar
}

// ParseTemplate parses a .env.template file where required vars have no
// default and optional vars are written as KEY=default.
// Lines starting with # are treated as comments attached to the next var.
func ParseTemplate(input string) (*Template, error) {
	tmpl := &Template{}
	lines := strings.Split(input, "\n")
	pendingComment := ""

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			pendingComment = ""
			continue
		}
		if strings.HasPrefix(line, "#") {
			pendingComment = strings.TrimSpace(strings.TrimPrefix(line, "#"))
			continue
		}
		idx := strings.Index(line, "=")
		if idx < 0 {
			return nil, fmt.Errorf("template: invalid line (missing '='): %q", line)
		}
		key := strings.TrimSpace(line[:idx])
		if key == "" {
			return nil, fmt.Errorf("template: empty key in line: %q", line)
		}
		val := strings.TrimSpace(line[idx+1:])
		v := TemplateVar{
			Key:      key,
			Default:  val,
			Required: val == "",
			Comment:  pendingComment,
		}
		tmpl.Vars = append(tmpl.Vars, v)
		pendingComment = ""
	}
	return tmpl, nil
}

// Validate checks that all required template vars are present in the provided
// env map. Returns a list of missing keys.
func (t *Template) Validate(env map[string]string) []string {
	var missing []string
	for _, v := range t.Vars {
		if v.Required {
			if val, ok := env[v.Key]; !ok || strings.TrimSpace(val) == "" {
				missing = append(missing, v.Key)
			}
		}
	}
	return missing
}

// Fill returns a new map seeded with template defaults, then overridden by env.
func (t *Template) Fill(env map[string]string) map[string]string {
	result := make(map[string]string, len(t.Vars))
	for _, v := range t.Vars {
		if !v.Required {
			result[v.Key] = v.Default
		}
	}
	for k, val := range env {
		result[k] = val
	}
	return result
}
