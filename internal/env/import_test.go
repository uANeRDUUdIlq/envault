package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempImportFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return p
}

func TestImportDotenv(t *testing.T) {
	p := writeTempImportFile(t, "FOO=bar\nBAZ=qux\n# comment\n\nEMPTY=\n")
	imp := NewImporter(ImportFormatDotenv, false)
	res, err := imp.ImportFile(p, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Imported != 3 {
		t.Errorf("expected 3 imported, got %d", res.Imported)
	}
	if res.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", res.Vars["FOO"])
	}
}

func TestImportJSON(t *testing.T) {
	p := writeTempImportFile(t, `{"KEY1":"val1","KEY2":"val2"}`)
	imp := NewImporter(ImportFormatJSON, false)
	res, err := imp.ImportFile(p, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Imported != 2 {
		t.Errorf("expected 2 imported, got %d", res.Imported)
	}
	if res.Vars["KEY1"] != "val1" {
		t.Errorf("expected KEY1=val1, got %q", res.Vars["KEY1"])
	}
}

func TestImportShell(t *testing.T) {
	p := writeTempImportFile(t, "export FOO=bar\nexport BAZ='hello world'\n")
	imp := NewImporter(ImportFormatShell, false)
	res, err := imp.ImportFile(p, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["BAZ"] != "hello world" {
		t.Errorf("expected BAZ='hello world', got %q", res.Vars["BAZ"])
	}
}

func TestImportSkipsExistingWithoutOverwrite(t *testing.T) {
	p := writeTempImportFile(t, "FOO=new\nBAR=new\n")
	imp := NewImporter(ImportFormatDotenv, false)
	res, err := imp.ImportFile(p, map[string]string{"FOO": "old"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "FOO" {
		t.Errorf("expected FOO skipped, got %v", res.Skipped)
	}
	if res.Imported != 1 {
		t.Errorf("expected 1 imported, got %d", res.Imported)
	}
}

func TestImportOverwriteExisting(t *testing.T) {
	p := writeTempImportFile(t, "FOO=new\n")
	imp := NewImporter(ImportFormatDotenv, true)
	res, err := imp.ImportFile(p, map[string]string{"FOO": "old"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["FOO"] != "new" {
		t.Errorf("expected FOO=new, got %q", res.Vars["FOO"])
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected no skipped, got %v", res.Skipped)
	}
}

func TestImportUnsupportedFormat(t *testing.T) {
	p := writeTempImportFile(t, "FOO=bar\n")
	imp := NewImporter(ImportFormat("toml"), false)
	_, err := imp.ImportFile(p, map[string]string{})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestImportFileNotFound(t *testing.T) {
	imp := NewImporter(ImportFormatDotenv, false)
	_, err := imp.ImportFile("/nonexistent/path.env", map[string]string{})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
