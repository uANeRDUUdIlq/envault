package env_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envault/envault/internal/env"
)

// TestExportRoundtrip verifies that exporting in dotenv format and re-parsing
// produces the original key-value pairs.
func TestExportRoundtrip(t *testing.T) {
	original := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"SECRET_KEY":  "s3cr3t!",
		"EMPTY_VALUE": "",
	}

	e := env.NewExporter(env.ExportOptions{Format: env.FormatDotenv})
	tmp := filepath.Join(t.TempDir(), "round.env")

	if err := e.Export(original, tmp); err != nil {
		t.Fatalf("Export: %v", err)
	}

	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	parsed, err := env.Parse(strings.NewReader(string(data)))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	for k, want := range original {
		got, ok := parsed[k]
		if !ok {
			t.Errorf("key %q missing after roundtrip", k)
			continue
		}
		if got != want {
			t.Errorf("key %q: want %q, got %q", k, want, got)
		}
	}
}

// TestExportKeyFilterIntegration ensures only requested keys appear in output.
func TestExportKeyFilterIntegration(t *testing.T) {
	vars := map[string]string{
		"PUBLIC":  "yes",
		"PRIVATE": "secret",
		"SHARED":  "common",
	}

	e := env.NewExporter(env.ExportOptions{
		Format: env.FormatDotenv,
		Keys:   []string{"PUBLIC", "SHARED"},
	})
	tmp := filepath.Join(t.TempDir(), "filtered.env")
	if err := e.Export(vars, tmp); err != nil {
		t.Fatalf("Export: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	content := string(data)

	if strings.Contains(content, "PRIVATE") {
		t.Errorf("PRIVATE should be excluded from export")
	}
	for _, k := range []string{"PUBLIC", "SHARED"} {
		if !strings.Contains(content, k) {
			t.Errorf("expected key %q in export output", k)
		}
	}
}
