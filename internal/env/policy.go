package env

import (
	"errors"
	"fmt"
	"strings"
)

// PolicyEffect defines whether a rule allows or denies access.
type PolicyEffect string

const (
	PolicyAllow PolicyEffect = "allow"
	PolicyDeny  PolicyEffect = "deny"
)

// PolicyRule describes a single rule in an environment policy.
type PolicyRule struct {
	Effect  PolicyEffect `json:"effect"`
	Keys    []string     `json:"keys"`
	Roles   []string     `json:"roles"`
	Comment string       `json:"comment,omitempty"`
}

// Policy holds a set of rules governing key access.
type Policy struct {
	Version string       `json:"version"`
	Rules   []PolicyRule `json:"rules"`
}

// EvaluationResult describes the outcome of a policy check.
type EvaluationResult struct {
	Allowed bool
	MatchedRule *PolicyRule
	Reason      string
}

// Evaluate checks whether the given role can access the given key.
// Rules are evaluated in order; the first match wins. Default is deny.
func (p *Policy) Evaluate(role, key string) EvaluationResult {
	for i := range p.Rules {
		rule := &p.Rules[i]
		if !ruleMatchesRole(rule, role) {
			continue
		}
		if !ruleMatchesKey(rule, key) {
			continue
		}
		allowed := rule.Effect == PolicyAllow
		reason := fmt.Sprintf("matched rule (effect=%s, comment=%q)", rule.Effect, rule.Comment)
		return EvaluationResult{Allowed: allowed, MatchedRule: rule, Reason: reason}
	}
	return EvaluationResult{Allowed: false, Reason: "no matching rule; default deny"}
}

func ruleMatchesRole(rule *PolicyRule, role string) bool {
	for _, r := range rule.Roles {
		if r == "*" || strings.EqualFold(r, role) {
			return true
		}
	}
	return false
}

func ruleMatchesKey(rule *PolicyRule, key string) bool {
	for _, k := range rule.Keys {
		if k == "*" || strings.EqualFold(k, key) {
			return true
		}
	}
	return false
}

// Validate checks that a Policy is structurally valid.
func (p *Policy) Validate() error {
	if p.Version == "" {
		return errors.New("policy: version is required")
	}
	for i, rule := range p.Rules {
		if rule.Effect != PolicyAllow && rule.Effect != PolicyDeny {
			return fmt.Errorf("policy: rule[%d] has invalid effect %q", i, rule.Effect)
		}
		if len(rule.Keys) == 0 {
			return fmt.Errorf("policy: rule[%d] must specify at least one key", i)
		}
		if len(rule.Roles) == 0 {
			return fmt.Errorf("policy: rule[%d] must specify at least one role", i)
		}
	}
	return nil
}
