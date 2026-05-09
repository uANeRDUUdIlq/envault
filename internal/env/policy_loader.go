package env

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const DefaultPolicyVersion = "1"

// LoadPolicy reads and parses a policy file from the given path.
func LoadPolicy(path string) (*Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("policy: file not found: %s", path)
		}
		return nil, fmt.Errorf("policy: read error: %w", err)
	}

	var p Policy
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("policy: invalid JSON: %w", err)
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	return &p, nil
}

// SavePolicy serialises a policy and writes it to the given path.
func SavePolicy(path string, p *Policy) error {
	if err := p.Validate(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("policy: marshal error: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("policy: write error: %w", err)
	}

	return nil
}

// CheckAccess is a convenience wrapper that loads a policy and evaluates it.
func CheckAccess(policyPath, role, key string) (EvaluationResult, error) {
	p, err := LoadPolicy(policyPath)
	if err != nil {
		return EvaluationResult{}, err
	}
	return p.Evaluate(role, key), nil
}
