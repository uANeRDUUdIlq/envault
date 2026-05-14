package env

import "fmt"

// InheritOptions controls how inheritance is applied.
type InheritOptions struct {
	// Overwrite existing keys in child with parent values.
	Overwrite bool
	// Keys to inherit; if empty, all keys are inherited.
	Keys []string
}

// InheritResult describes the outcome of an inherit operation.
type InheritResult struct {
	Inherited []string
	Skipped   []string
}

// SummaryString returns a human-readable summary of the result.
func (r InheritResult) SummaryString() string {
	return fmt.Sprintf("inherited %d key(s), skipped %d key(s)",
		len(r.Inherited), len(r.Skipped))
}

// Inherit copies keys from parent into child according to opts.
// Parent and child are maps of env var name -> value.
func Inherit(parent, child map[string]string, opts InheritOptions) (map[string]string, InheritResult) {
	result := map[string]string{}
	for k, v := range child {
		result[k] = v
	}

	var res InheritResult

	wantKey := func(k string) bool {
		if len(opts.Keys) == 0 {
			return true
		}
		for _, wanted := range opts.Keys {
			if wanted == k {
				return true
			}
		}
		return false
	}

	for k, v := range parent {
		if !wantKey(k) {
			continue
		}
		if _, exists := result[k]; exists && !opts.Overwrite {
			res.Skipped = append(res.Skipped, k)
			continue
		}
		result[k] = v
		res.Inherited = append(res.Inherited, k)
	}

	return result, res
}
