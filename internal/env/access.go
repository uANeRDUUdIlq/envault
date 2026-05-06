package env

import (
	"errors"
	"sort"
	"strings"
)

// AccessPolicy defines which keys a given role or user may read or write.
type AccessPolicy struct {
	Role      string
	AllowRead []string // glob-style prefixes; "*" means all
	AllowWrite []string
}

// AccessStore holds a collection of policies keyed by role name.
type AccessStore struct {
	policies map[string]AccessPolicy
}

// NewAccessStore returns an empty AccessStore.
func NewAccessStore() *AccessStore {
	return &AccessStore{policies: make(map[string]AccessPolicy)}
}

// SetPolicy registers or replaces the policy for the given role.
func (a *AccessStore) SetPolicy(p AccessPolicy) error {
	if strings.TrimSpace(p.Role) == "" {
		return errors.New("access: role name must not be empty")
	}
	a.policies[p.Role] = p
	return nil
}

// GetPolicy returns the policy for a role, or false if not found.
func (a *AccessStore) GetPolicy(role string) (AccessPolicy, bool) {
	p, ok := a.policies[role]
	return p, ok
}

// DeletePolicy removes the policy for the given role.
func (a *AccessStore) DeletePolicy(role string) {
	delete(a.policies, role)
}

// Roles returns all registered role names in sorted order.
func (a *AccessStore) Roles() []string {
	names := make([]string, 0, len(a.policies))
	for k := range a.policies {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// CanRead reports whether the given role may read the supplied key.
func (a *AccessStore) CanRead(role, key string) bool {
	p, ok := a.policies[role]
	if !ok {
		return false
	}
	return matchesAny(p.AllowRead, key)
}

// CanWrite reports whether the given role may write the supplied key.
func (a *AccessStore) CanWrite(role, key string) bool {
	p, ok := a.policies[role]
	if !ok {
		return false
	}
	return matchesAny(p.AllowWrite, key)
}

// FilterReadable returns only the key/value pairs the role is allowed to read.
func (a *AccessStore) FilterReadable(role string, vars map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range vars {
		if a.CanRead(role, k) {
			out[k] = v
		}
	}
	return out
}

// matchesAny returns true if key matches any prefix in the list, or the list
// contains the wildcard "*".
func matchesAny(patterns []string, key string) bool {
	for _, p := range patterns {
		if p == "*" {
			return true
		}
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}
