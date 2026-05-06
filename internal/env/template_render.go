package env

import (
	"fmt"
	"strings"
)

// RenderOptions controls how a template is rendered to text.
type RenderOptions struct {
	IncludeComments bool
	MaskSecrets     bool
	SecretSuffixes  []string
}

// Render produces a .env-style string from the template, optionally filling
// in values from the provided env map.
func (t *Template) Render(env map[string]string, opts RenderOptions) string {
	var sb strings.Builder
	for _, v := range t.Vars {
		if opts.IncludeComments && v.Comment != "" {
			fmt.Fprintf(&sb, "# %s\n", v.Comment)
		}
		val := v.Default
		if override, ok := env[v.Key]; ok {
			val = override
		}
		if opts.MaskSecrets && isSecret(v.Key, opts.SecretSuffixes) && val != "" {
			val = "***"
		}
		fmt.Fprintf(&sb, "%s=%s\n", v.Key, val)
	}
	return sb.String()
}

// isSecret returns true if the key ends with one of the given suffixes
// (case-insensitive). Falls back to common secret keywords when no suffixes
// are provided.
func isSecret(key string, suffixes []string) bool {
	upper := strings.ToUpper(key)
	if len(suffixes) == 0 {
		suffixes = []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PASS"}
	}
	for _, s := range suffixes {
		if strings.HasSuffix(upper, strings.ToUpper(s)) {
			return true
		}
	}
	return false
}

// MissingKeys returns template vars that are required but absent from env.
func (t *Template) MissingKeys(env map[string]string) []string {
	return t.Validate(env)
}

// ExtraKeys returns keys present in env that are not declared in the template.
func (t *Template) ExtraKeys(env map[string]string) []string {
	declared := make(map[string]bool, len(t.Vars))
	for _, v := range t.Vars {
		declared[v.Key] = true
	}
	var extra []string
	for k := range env {
		if !declared[k] {
			extra = append(extra, k)
		}
	}
	return extra
}
