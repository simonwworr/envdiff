package reconcile

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeExportEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "APP_NAME", OldValue: "myapp", NewValue: "myapp", Status: diff.StatusUnchanged},
		{Key: "SECRET_KEY", OldValue: "old", NewValue: "new", Status: diff.StatusChanged},
		{Key: "ADDED_VAR", OldValue: "", NewValue: "hello world", Status: diff.StatusAdded},
	}
}

func TestExport_DotenvFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, makeExportEntries(), ExportOptions{Format: ExportFormatDotenv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("expected APP_NAME line, got:\n%s", out)
	}
	if !strings.Contains(out, `ADDED_VAR="hello world"`) {
		t.Errorf("expected quoted ADDED_VAR, got:\n%s", out)
	}
}

func TestExport_ShellFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, makeExportEntries(), ExportOptions{Format: ExportFormatShell})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export APP_NAME=") {
		t.Errorf("expected shell export prefix, got:\n%s", out)
	}
}

func TestExport_DockerFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, makeExportEntries(), ExportOptions{Format: ExportFormatDocker})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "--env APP_NAME=") {
		t.Errorf("expected docker --env prefix, got:\n%s", out)
	}
}

func TestExport_OnlyChanged(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, makeExportEntries(), ExportOptions{
		Format:      ExportFormatDotenv,
		OnlyChanged: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "APP_NAME") {
		t.Errorf("unchanged key should be excluded, got:\n%s", out)
	}
	if !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("changed key should be included, got:\n%s", out)
	}
}

func TestExport_MasksSensitiveKeys(t *testing.T) {
	masker := diff.NewMasker()
	var buf bytes.Buffer
	err := Export(&buf, makeExportEntries(), ExportOptions{
		Format: ExportFormatDotenv,
		Masker: masker,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "new") && strings.Contains(out, "SECRET_KEY") {
		t.Errorf("sensitive key value should be masked, got:\n%s", out)
	}
}

func TestExportToFile_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "exported.env")
	err := ExportToFile(path, makeExportEntries(), ExportOptions{Format: ExportFormatDotenv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read exported file: %v", err)
	}
	if !strings.Contains(string(data), "APP_NAME") {
		t.Errorf("expected APP_NAME in exported file, got:\n%s", data)
	}
}
