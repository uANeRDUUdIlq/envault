package env_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/envault/envault/internal/env"
)

func writeIntegrationFile(t *testing.T, name, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	return p
}

func TestImportExportRoundtrip(t *testing.T) {
	vars := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "production",
	}

	// Serialize to dotenv format.
	dotenvContent := env.Serialize(vars)
	p := writeIntegrationFile(t, "roundtrip.env", dotenvContent)

	imp := env.NewImporter(env.ImportFormatDotenv, false)
	res, err := imp.ImportFile(p, map[string]string{})
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	if res.Imported != len(vars) {
		t.Errorf("expected %d imported, got %d", len(vars), res.Imported)
	}
	for k, v := range vars {
		if res.Vars[k] != v {
			t.Errorf("key %s: expected %q, got %q", k, v, res.Vars[k])
		}
	}
}

func TestImportJSONRoundtrip(t *testing.T) {
	vars := map[string]string{
		"SECRET_KEY": "abc123",
		"API_URL":    "https://api.example.com",
	}

	data, err := json.Marshal(vars)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	p := writeIntegrationFile(t, "vars.json", string(data))

	imp := env.NewImporter(env.ImportFormatJSON, false)
	res, err := imp.ImportFile(p, map[string]string{})
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	if res.Vars["SECRET_KEY"] != "abc123" {
		t.Errorf("expected SECRET_KEY=abc123, got %q", res.Vars["SECRET_KEY"])
	}
}
