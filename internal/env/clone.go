package env

import "fmt"

// CloneResult holds the outcome of a clone operation.
type CloneResult struct {
	SourceNamespace string
	DestNamespace   string
	Copied          int
	Skipped         int
	Overwritten     int
}

// SummaryString returns a human-readable summary of the clone result.
func (r CloneResult) SummaryString() string {
	return fmt.Sprintf(
		"cloned %s -> %s: %d copied, %d skipped, %d overwritten",
		r.SourceNamespace, r.DestNamespace, r.Copied, r.Skipped, r.Overwritten,
	)
}

// CloneOptions controls how a clone operation behaves.
type CloneOptions struct {
	Overwrite bool
	FilterKeys []string // if non-empty, only clone these keys
}

// Cloner copies all (or a filtered subset of) env vars from one namespace to another.
type Cloner struct {
	src  map[string]string
	dest map[string]string
	opts CloneOptions
}

// NewCloner creates a Cloner that will copy from src into dest.
func NewCloner(src, dest map[string]string, opts CloneOptions) *Cloner {
	return &Cloner{src: src, dest: dest, opts: opts}
}

// Clone performs the copy and returns the result along with the updated dest map.
func (c *Cloner) Clone(srcNS, destNS string) (map[string]string, CloneResult) {
	result := CloneResult{
		SourceNamespace: srcNS,
		DestNamespace:   destNS,
	}

	filter := make(map[string]bool, len(c.opts.FilterKeys))
	for _, k := range c.opts.FilterKeys {
		filter[k] = true
	}

	out := make(map[string]string, len(c.dest))
	for k, v := range c.dest {
		out[k] = v
	}

	for k, v := range c.src {
		if len(filter) > 0 && !filter[k] {
			continue
		}
		if existing, exists := out[k]; exists {
			if !c.opts.Overwrite {
				result.Skipped++
				continue
			}
			if existing != v {
				result.Overwritten++
			}
		} else {
			result.Copied++
		}
		out[k] = v
	}

	return out, result
}
