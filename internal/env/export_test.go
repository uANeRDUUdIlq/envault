package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExportDotenv(t *testing.T) {
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	e := NewExporter(ExportOptions{Format: FormatDotenv})

	tmp := filepath.Join(t.TempDir(), "out.env")
	if err := e.Export(vars, tmp); err != nil {
		t.Fatalf("Export: %v", err)
	}

	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	content := string(data)
	for k, v := range vars {
		expected := k + "=" + v
		if !strings.Contains(content, expected) {
			t.Errorf("expected %q in output, got:\n%s", expected, content)
		}
	}
}

func TestExportShellNoExport(t *testing.T) {
	vars := map[string]string{"KEY": "value with spaces"}
	e := NewExporter(ExportOptions{Format: FormatShell, Export: false})

	tmp := filepath.Join(t.TempDir(), "out.sh")
	if err := e.Export(vars, tmp); err != nil {
		t.Fatalf("Export: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	if strings.Contains(string(data), "export ") {
		t.Errorf("did not expect 'export' prefix, got: %s", data)
	}
}

func TestExportShellWithExport(t *testing.T) {
	vars := map[string]string{"KEY": "val"}
	e := NewExporter(ExportOptions{Format: FormatShell, Export: true})

	tmp := filepath.Join(t.TempDir(), "out.sh")
	if err := e.Export(vars, tmp); err != nil {
		t.Fatalf("Export: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	if !strings.HasPrefix(string(data), "export ") {
		t.Errorf("expected 'export' prefix, got: %s", data)
	}
}

func TestExportJSON(t *testing.T) {
	vars := map[string]string{"A": "1"}
	e := NewExporter(ExportOptions{Format: FormatJSON})

	tmp := filepath.Join(t.TempDir(), "out.json")
	if err := e.Export(vars, tmp); err != nil {
		t.Fatalf("Export: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	if !strings.Contains(string(data), "\"A\"") {
		t.Errorf("expected JSON key A, got: %s", data)
	}
}

func TestExportFilterKeys(t *testing.T) {
	vars := map[string]string{"KEEP": "yes", "SKIP": "no"}
	e := NewExporter(ExportOptions{Format: FormatDotenv, Keys: []string{"KEEP"}})

	tmp := filepath.Join(t.TempDir(), "filtered.env")
	if err := e.Export(vars, tmp); err != nil {
		t.Fatalf("Export: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	content := string(data)
	if strings.Contains(content, "SKIP") {
		t.Errorf("expected SKIP to be filtered out, got: %s", content)
	}
	if !strings.Contains(content, "KEEP") {
		t.Errorf("expected KEEP in output, got: %s", content)
	}
}

func TestExportDefaultFormat(t *testing.T) {
	e := NewExporter(ExportOptions{})
	if e.opts.Format != FormatDotenv {
		t.Errorf("expected default format dotenv, got %q", e.opts.Format)
	}
}
