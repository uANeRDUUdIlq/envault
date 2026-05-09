package env

import (
	"testing"
)

// TestPolicyRoundtripAndEvaluate saves a policy to disk, loads it back,
// and verifies evaluation is consistent across the full cycle.
func TestPolicyRoundtripAndEvaluate(t *testing.T) {
	orig := &Policy{
		Version: "1",
		Rules: []PolicyRule{
			{Effect: PolicyAllow, Keys: []string{"SECRET"}, Roles: []string{"ops"}, Comment: "ops only"},
			{Effect: PolicyDeny, Keys: []string{"*"}, Roles: []string{"*"}, Comment: "deny all others"},
		},
	}

	path := writePolicyFile(t, orig)

	cases := []struct {
		role    string
		key     string
		wantOK  bool
	}{
		{"ops", "SECRET", true},
		{"ops", "OTHER", false},
		{"dev", "SECRET", false},
		{"dev", "OTHER", false},
	}

	for _, tc := range cases {
		res, err := CheckAccess(path, tc.role, tc.key)
		if err != nil {
			t.Fatalf("CheckAccess(%s,%s): %v", tc.role, tc.key, err)
		}
		if res.Allowed != tc.wantOK {
			t.Errorf("CheckAccess(%s,%s): got allowed=%v, want %v (reason: %s)",
				tc.role, tc.key, res.Allowed, tc.wantOK, res.Reason)
		}
	}
}

// TestPolicyWildcardRoleAndKey verifies that wildcard matching works for both
// roles and keys simultaneously.
func TestPolicyWildcardRoleAndKey(t *testing.T) {
	p := &Policy{
		Version: "1",
		Rules: []PolicyRule{
			{Effect: PolicyAllow, Keys: []string{"*"}, Roles: []string{"*"}, Comment: "allow all"},
		},
	}

	for _, role := range []string{"admin", "dev", "guest"} {
		for _, key := range []string{"FOO", "BAR", "SECRET"} {
			res := p.Evaluate(role, key)
			if !res.Allowed {
				t.Errorf("expected allow for role=%s key=%s, got deny", role, key)
			}
		}
	}
}
