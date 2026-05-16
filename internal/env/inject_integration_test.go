package env

import (
	"os"
	"testing"
)

// TestInjectRoundtripWithParser parses a dotenv string and injects it,
// then ejects and verifies the keys are gone.
func TestInjectRoundtripWithParser(t *testing.T) {
	input := "DB_HOST=localhost\nDB_PORT=5432\n"
	vars, err := Parse(input)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	for k := range vars {
		os.Unsetenv(k)
	}
	defer func() {
		for k := range vars {
			os.Unsetenv(k)
		}
	}()

	inj := NewInjector(InjectOptions{Overwrite: true})
	res, err := inj.Inject(vars)
	if err != nil {
		t.Fatalf("inject: %v", err)
	}
	if len(res.Injected) != 2 {
		t.Errorf("expected 2 injected, got %d", len(res.Injected))
	}

	if os.Getenv("DB_HOST") != "localhost" {
		t.Errorf("DB_HOST not set correctly")
	}

	removed := inj.Eject(vars)
	if len(removed) != 2 {
		t.Errorf("expected 2 ejected, got %d", len(removed))
	}
	if _, ok := os.LookupEnv("DB_HOST"); ok {
		t.Errorf("DB_HOST should have been ejected")
	}
}

// TestInjectPrefixAndFilterIntegration tests prefix + key filtering together.
func TestInjectPrefixAndFilterIntegration(t *testing.T) {
	input := "SECRET_KEY=abc\nPUBLIC_URL=https://example.com\nDEBUG=true\n"
	vars, err := Parse(input)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	defer func() {
		for k := range vars {
			os.Unsetenv("ENV_" + k)
		}
	}()

	inj := NewInjector(InjectOptions{
		Overwrite: true,
		Prefix:    "ENV_",
		Keys:      []string{"SECRET_KEY", "PUBLIC_URL"},
	})
	res, err := inj.Inject(vars)
	if err != nil {
		t.Fatalf("inject: %v", err)
	}
	if len(res.Injected) != 2 {
		t.Errorf("expected 2 injected, got %d", len(res.Injected))
	}
	if _, ok := os.LookupEnv("ENV_DEBUG"); ok {
		t.Errorf("ENV_DEBUG should not have been injected")
	}
	if os.Getenv("ENV_SECRET_KEY") != "abc" {
		t.Errorf("ENV_SECRET_KEY not set correctly")
	}
}
