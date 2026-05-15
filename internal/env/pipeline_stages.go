package env

import (
	"fmt"
	"strings"
)

// PrefixStage returns a PipelineStage that adds a prefix to all keys.
func PrefixStage(prefix string) PipelineStage {
	return PipelineStage{
		Name: fmt.Sprintf("prefix(%s)", prefix),
		Transform: func(vars map[string]string) (map[string]string, error) {
			out := make(map[string]string, len(vars))
			for k, v := range vars {
				out[prefix+k] = v
			}
			return out, nil
		},
	}
}

// UppercaseKeysStage returns a PipelineStage that uppercases all keys.
func UppercaseKeysStage() PipelineStage {
	return PipelineStage{
		Name: "uppercase_keys",
		Transform: func(vars map[string]string) (map[string]string, error) {
			out := make(map[string]string, len(vars))
			for k, v := range vars {
				out[strings.ToUpper(k)] = v
			}
			return out, nil
		},
	}
}

// FilterKeysStage returns a PipelineStage that retains only the given keys.
func FilterKeysStage(keys []string) PipelineStage {
	allowed := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		allowed[k] = struct{}{}
	}
	return PipelineStage{
		Name: "filter_keys",
		Transform: func(vars map[string]string) (map[string]string, error) {
			out := make(map[string]string)
			for k, v := range vars {
				if _, ok := allowed[k]; ok {
					out[k] = v
				}
			}
			return out, nil
		},
	}
}

// RequireKeysStage returns a PipelineStage that errors if any required key is missing.
func RequireKeysStage(keys []string) PipelineStage {
	return PipelineStage{
		Name: "require_keys",
		Transform: func(vars map[string]string) (map[string]string, error) {
			for _, k := range keys {
				if _, ok := vars[k]; !ok {
					return nil, fmt.Errorf("required key %q is missing", k)
				}
			}
			return copyVars(vars), nil
		},
	}
}
