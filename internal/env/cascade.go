package env

import "fmt"

// CascadeOptions controls how cascading merges behave.
type CascadeOptions struct {
	// Overwrite allows later layers to overwrite earlier ones.
	Overwrite bool
	// Layers is an ordered list of environment names (e.g. ["base", "staging", "prod"]).
	Layers []string
}

// CascadeResult holds the merged variables and provenance info.
type CascadeResult struct {
	Vars     map[string]string
	// Origin maps each key to the layer name it came from.
	Origin   map[string]string
}

// SummaryLines returns a human-readable summary of the cascade result.
func (r *CascadeResult) SummaryLines() []string {
	lines := make([]string, 0, len(r.Vars))
	for _, k := range sortedKeys(r.Vars) {
		lines = append(lines, fmt.Sprintf("%s (from %s)", k, r.Origin[k]))
	}
	return lines
}

// Cascade merges multiple layers of env vars in order.
// layers is a map from layer name to its variables.
// The order of application is defined by opts.Layers.
func Cascade(layers map[string]map[string]string, opts CascadeOptions) (*CascadeResult, error) {
	result := &CascadeResult{
		Vars:   make(map[string]string),
		Origin: make(map[string]string),
	}

	for _, name := range opts.Layers {
		vars, ok := layers[name]
		if !ok {
			return nil, fmt.Errorf("cascade: layer %q not found", name)
		}
		for k, v := range vars {
			if _, exists := result.Vars[k]; exists && !opts.Overwrite {
				continue
			}
			result.Vars[k] = v
			result.Origin[k] = name
		}
	}

	return result, nil
}
