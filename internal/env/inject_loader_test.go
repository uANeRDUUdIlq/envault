package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempInjectFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("write temp inject file: %v", err)
	}
	return p
}

func TestInjectFromFile(t *testing.T) {
	p := writeTempInjectFile(t, "LOADER_A=foo\nLOADER_B=bar\n")
	defer os.Unsetenv("LOADER_A")
	defer os.Unsetenv("LOADER_B")

	res, err := InjectFromFile(p, InjectOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("InjectFromFile: %v", err)
	}
	if len(res.Injected) != 2 {
		t.Errorf("expected 2 injected, got %d", len(res.Injected))
	}
	if os.Getenv("LOADER_A") != "foo" {
		t.Errorf("LOADER_A not set")
	}
}

func TestInjectFromFileMissing(t *testing.T) {
	_, err := InjectFromFile("/nonexistent/.env", InjectOptions{})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestEjectFromFile(t *testing.T) {
	p := writeTempInjectFile(t, "EJECT_LOADER_X=1\n")
	os.Setenv("EJECT_LOADER_X", "1")

	removed, err := EjectFromFile(p, InjectOptions{})
	if err != nil {
		t.Fatalf("EjectFromFile: %v", err)
	}
	if len(removed) != 1 {
		t.Errorf("expected 1 removed, got %d", len(removed))
	}
	if _, ok := os.LookupEnv("EJECT_LOADER_X"); ok {
		t.Errorf("EJECT_LOADER_X should have been removed")
	}
}

func TestInjectFromFileWithPrefix(t *testing.T) {
	p := writeTempInjectFile(t, "HOST=db\nPORT=5432\n")
	defer os.Unsetenv("PFX_HOST")
	defer os.Unsetenv("PFX_PORT")

	res, err := InjectFromFile(p, InjectOptions{Overwrite: true, Prefix: "PFX_"})
	if err != nil {
		t.Fatalf("InjectFromFile with prefix: %v", err)
	}
	if len(res.Injected) != 2 {
		t.Errorf("expected 2 injected, got %d", len(res.Injected))
	}
	if os.Getenv("PFX_HOST") != "db" {
		t.Errorf("PFX_HOST not set correctly")
	}
}
