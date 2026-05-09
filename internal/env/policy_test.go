package env

import (
	"testing"
)

func samplePolicy() *Policy {
	return &Policy{
		Version: "1",
		Rules: []PolicyRule{
			{
				Effect:  PolicyAllow,
				Keys:    []string{"DB_URL", "DB_PASS"},
				Roles:   []string{"admin"},
				Comment: "admins can read DB creds",
			},
			{
				Effect: PolicyDeny,
				Keys:   []string{"DB_PASS"},
				Roles:  []string{"developer"},
			},
			{
				Effect: PolicyAllow,
				Keys:   []string{"*"},
				Roles:  []string{"developer"},
			},
		},
	}
}

func TestPolicyEvaluateAllow(t *testing.T) {
	p := samplePolicy()
	res := p.Evaluate("admin", "DB_URL")
	if !res.Allowed {
		t.Fatalf("expected allow, got deny: %s", res.Reason)
	}
}

func TestPolicyEvaluateDeny(t *testing.T) {
	p := samplePolicy()
	res := p.Evaluate("developer", "DB_PASS")
	if res.Allowed {
		t.Fatalf("expected deny, got allow: %s", res.Reason)
	}
}

func TestPolicyEvaluateWildcardKey(t *testing.T) {
	p := samplePolicy()
	res := p.Evaluate("developer", "API_KEY")
	if !res.Allowed {
		t.Fatalf("expected allow via wildcard, got deny: %s", res.Reason)
	}
}

func TestPolicyEvaluateNoMatchDefaultDeny(t *testing.T) {
	p := samplePolicy()
	res := p.Evaluate("stranger", "DB_URL")
	if res.Allowed {
		t.Fatalf("expected default deny, got allow")
	}
}

func TestPolicyValidateMissingVersion(t *testing.T) {
	p := &Policy{Rules: []PolicyRule{{Effect: PolicyAllow, Keys: []string{"K"}, Roles: []string{"r"}}}}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for missing version")
	}
}

func TestPolicyValidateInvalidEffect(t *testing.T) {
	p := &Policy{Version: "1", Rules: []PolicyRule{{Effect: "maybe", Keys: []string{"K"}, Roles: []string{"r"}}}}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for invalid effect")
	}
}

func TestPolicyValidateEmptyKeys(t *testing.T) {
	p := &Policy{Version: "1", Rules: []PolicyRule{{Effect: PolicyAllow, Keys: []string{}, Roles: []string{"r"}}}}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for empty keys")
	}
}

func TestPolicyValidateEmptyRoles(t *testing.T) {
	p := &Policy{Version: "1", Rules: []PolicyRule{{Effect: PolicyAllow, Keys: []string{"K"}, Roles: []string{}}}}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for empty roles")
	}
}
