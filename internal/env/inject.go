package env

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// InjectOptions controls how variables are injected into the process environment.
type InjectOptions struct {
	// Overwrite existing environment variables if true.
	Overwrite bool
	// Prefix prepended to every injected key.
	Prefix string
	// Keys restricts injection to only these keys (empty means all).
	Keys []string
}

// InjectResult summarises what happened during an injection.
type InjectResult struct {
	Injected []string
	Skipped  []string
	Prefixed bool
}

// SummaryString returns a human-readable summary of the injection.
func (r *InjectResult) SummaryString() string {
	return fmt.Sprintf("injected=%d skipped=%d", len(r.Injected), len(r.Skipped))
}

// Injector writes env vars into the current process environment.
type Injector struct {
	opts InjectOptions
}

// NewInjector creates an Injector with the given options.
func NewInjector(opts InjectOptions) *Injector {
	return &Injector{opts: opts}
}

// Inject applies vars to os.Environ, respecting options.
func (inj *Injector) Inject(vars map[string]string) (*InjectResult, error) {
	allowed := allowedSet(inj.opts.Keys)
	result := &InjectResult{Prefixed: inj.opts.Prefix != ""}

	keys := sortedKeys(vars)
	for _, k := range keys {
		if len(allowed) > 0 && !allowed[k] {
			continue
		}
		env := strings.ToUpper(inj.opts.Prefix + k)
		_, exists := os.LookupEnv(env)
		if exists && !inj.opts.Overwrite {
			result.Skipped = append(result.Skipped, env)
			continue
		}
		if err := os.Setenv(env, vars[k]); err != nil {
			return nil, fmt.Errorf("inject: setenv %s: %w", env, err)
		}
		result.Injected = append(result.Injected, env)
	}
	return result, nil
}

// Eject removes previously injected vars from the process environment.
func (inj *Injector) Eject(vars map[string]string) []string {
	var removed []string
	for k := range vars {
		env := strings.ToUpper(inj.opts.Prefix + k)
		if _, ok := os.LookupEnv(env); ok {
			os.Unsetenv(env)
			removed = append(removed, env)
		}
	}
	sort.Strings(removed)
	return removed
}

func allowedSet(keys []string) map[string]bool {
	if len(keys) == 0 {
		return nil
	}
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}
