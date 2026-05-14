package env

import (
	"regexp"
	"strings"
)

// RedactOptions controls how values are redacted.
type RedactOptions struct {
	// MaskChar is the character used to replace secret values.
	MaskChar string
	// PartialReveal shows the first N characters before masking.
	PartialReveal int
	// Keys is an explicit list of keys to redact; if empty, heuristics apply.
	Keys []string
}

// DefaultRedactOptions returns sensible defaults.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{
		MaskChar:      "*",
		PartialReveal: 0,
	}
}

var secretKeyPattern = regexp.MustCompile(
	`(?i)(password|secret|token|key|auth|credential|private|api_?key|access_?key|passwd)`,
)

// Redactor masks sensitive values in a map of env vars.
type Redactor struct {
	opts RedactOptions
	explicit map[string]struct{}
}

// NewRedactor creates a Redactor with the given options.
func NewRedactor(opts RedactOptions) *Redactor {
	explicit := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		explicit[strings.ToUpper(k)] = struct{}{}
	}
	return &Redactor{opts: opts, explicit: explicit}
}

// Redact returns a copy of vars with sensitive values masked.
func (r *Redactor) Redact(vars map[string]string) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		if r.shouldRedact(k) {
			out[k] = r.mask(v)
		} else {
			out[k] = v
		}
	}
	return out
}

// RedactedKeys returns the list of keys that would be redacted.
func (r *Redactor) RedactedKeys(vars map[string]string) []string {
	var keys []string
	for k := range vars {
		if r.shouldRedact(k) {
			keys = append(keys, k)
		}
	}
	return keys
}

func (r *Redactor) shouldRedact(key string) bool {
	if _, ok := r.explicit[strings.ToUpper(key)]; ok {
		return true
	}
	return secretKeyPattern.MatchString(key)
}

func (r *Redactor) mask(value string) string {
	if value == "" {
		return ""
	}
	maskLen := len(value) - r.opts.PartialReveal
	if maskLen <= 0 {
		return value
	}
	prefix := ""
	if r.opts.PartialReveal > 0 && r.opts.PartialReveal < len(value) {
		prefix = value[:r.opts.PartialReveal]
	}
	return prefix + strings.Repeat(r.opts.MaskChar, maskLen)
}
