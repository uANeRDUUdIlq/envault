package env

import "testing"

// TestInheritRoundtripWithParser verifies that Inherit works correctly
// when parent and child are produced by the env parser.
func TestInheritRoundtripWithParser(t *testing.T) {
	parentSrc := "BASE_URL=https://example.com\nTIMEOUT=30\nDEBUG=false\n"
	childSrc := "APP_NAME=myapp\nTIMEOUT=60\n"

	parentVars, err := Parse(parentSrc)
	if err != nil {
		t.Fatalf("parse parent: %v", err)
	}
	childVars, err := Parse(childSrc)
	if err != nil {
		t.Fatalf("parse child: %v", err)
	}

	out, res := Inherit(parentVars, childVars, InheritOptions{Overwrite: false})

	if out["BASE_URL"] != "https://example.com" {
		t.Errorf("BASE_URL not inherited: %s", out["BASE_URL"])
	}
	if out["TIMEOUT"] != "60" {
		t.Errorf("child TIMEOUT should be preserved: %s", out["TIMEOUT"])
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should remain: %s", out["APP_NAME"])
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "TIMEOUT" {
		t.Errorf("expected TIMEOUT skipped, got %v", res.Skipped)
	}

	// Serialize and re-parse to ensure output is valid.
	serialised := Serialize(out)
	reparsed, err := Parse(serialised)
	if err != nil {
		t.Fatalf("reparse: %v", err)
	}
	if reparsed["BASE_URL"] != "https://example.com" {
		t.Errorf("BASE_URL lost after serialise/parse: %s", reparsed["BASE_URL"])
	}
}

// TestInheritKeyFilterIntegration verifies selective key inheritance end-to-end.
func TestInheritKeyFilterIntegration(t *testing.T) {
	parentSrc := "DB_HOST=db.prod\nDB_PORT=5432\nSECRET_KEY=s3cr3t\n"
	childSrc := "APP=web\n"

	parentVars, _ := Parse(parentSrc)
	childVars, _ := Parse(childSrc)

	out, res := Inherit(parentVars, childVars, InheritOptions{
		Keys: []string{"DB_HOST", "DB_PORT"},
	})

	if _, ok := out["SECRET_KEY"]; ok {
		t.Error("SECRET_KEY should not be inherited")
	}
	if out["DB_HOST"] != "db.prod" {
		t.Errorf("DB_HOST not inherited: %s", out["DB_HOST"])
	}
	if len(res.Inherited) != 2 {
		t.Errorf("expected 2 inherited, got %d", len(res.Inherited))
	}
}
