package env

import "testing"

// TestCloneRoundtripWithParser verifies that parsed env vars can be cloned
// and the result serializes correctly.
func TestCloneRoundtripWithParser(t *testing.T) {
	raw := "DB_HOST=localhost\nDB_PORT=5432\nSECRET_KEY=abc123\n"
	vars, err := Parse(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	dest := map[string]string{"APP_ENV": "production"}
	c := NewCloner(vars, dest, CloneOptions{Overwrite: false})
	out, res := c.Clone("dev", "prod")

	if res.Copied != 3 {
		t.Fatalf("expected 3 copied, got %d", res.Copied)
	}
	if out["APP_ENV"] != "production" {
		t.Fatal("expected APP_ENV to be preserved in dest")
	}

	serialised := Serialize(out)
	if len(serialised) == 0 {
		t.Fatal("expected non-empty serialized output")
	}
}

// TestCloneFilterAndOverwriteIntegration tests combined filter + overwrite logic.
func TestCloneFilterAndOverwriteIntegration(t *testing.T) {
	src := map[string]string{
		"DB_HOST": "new-host",
		"DB_PORT": "5433",
		"DEBUG":   "true",
	}
	dest := map[string]string{
		"DB_HOST": "old-host",
		"DB_PORT": "5432",
	}

	opts := CloneOptions{
		Overwrite:  true,
		FilterKeys: []string{"DB_HOST"},
	}
	c := NewCloner(src, dest, opts)
	out, res := c.Clone("staging", "prod")

	if res.Overwritten != 1 {
		t.Fatalf("expected 1 overwritten, got %d", res.Overwritten)
	}
	if out["DB_HOST"] != "new-host" {
		t.Fatalf("expected DB_HOST=new-host, got %s", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5432" {
		t.Fatalf("expected DB_PORT unchanged, got %s", out["DB_PORT"])
	}
	if _, ok := out["DEBUG"]; ok {
		t.Fatal("DEBUG should not be in output due to filter")
	}
}
