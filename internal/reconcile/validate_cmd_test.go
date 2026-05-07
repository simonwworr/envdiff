package reconcile

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/config"
)

func writeTempValidateEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return p
}

func TestRunValidate_PassesWhenNoIssues(t *testing.T) {
	dir := t.TempDir()
	base := writeTempValidateEnv(t, dir, "base.env", "HOST=localhost\nPORT=8080\n")
	other := writeTempValidateEnv(t, dir, "other.env", "HOST=localhost\nPORT=8080\n")

	cfg := &config.Config{BaseFile: base, OtherFile: other}
	var buf bytes.Buffer
	code := RunValidate(cfg, false, &buf)
	if code != 0 {
		t.Errorf("expected exit code 0, got %d; output: %s", code, buf.String())
	}
	if !stringContains(buf.String(), "validation passed") {
		t.Errorf("expected success message, got: %s", buf.String())
	}
}

func TestRunValidate_FailsWithRequireAll(t *testing.T) {
	dir := t.TempDir()
	base := writeTempValidateEnv(t, dir, "base.env", "HOST=localhost\n")
	other := writeTempValidateEnv(t, dir, "other.env", "HOST=localhost\nNEW_VAR=added\n")

	cfg := &config.Config{BaseFile: base, OtherFile: other}
	var buf bytes.Buffer
	code := RunValidate(cfg, true, &buf)
	if code != 2 {
		t.Errorf("expected exit code 2, got %d; output: %s", code, buf.String())
	}
	if !stringContains(buf.String(), "NEW_VAR") {
		t.Errorf("expected NEW_VAR in output, got: %s", buf.String())
	}
}

func TestRunValidate_MissingBaseFile(t *testing.T) {
	dir := t.TempDir()
	other := writeTempValidateEnv(t, dir, "other.env", "KEY=val\n")

	cfg := &config.Config{BaseFile: filepath.Join(dir, "nonexistent.env"), OtherFile: other}
	var buf bytes.Buffer
	code := RunValidate(cfg, false, &buf)
	if code != 1 {
		t.Errorf("expected exit code 1 for missing base file, got %d", code)
	}
}

func TestRunValidate_ChangedEmptyBase(t *testing.T) {
	dir := t.TempDir()
	base := writeTempValidateEnv(t, dir, "base.env", "SECRET=\n")
	other := writeTempValidateEnv(t, dir, "other.env", "SECRET=newvalue\n")

	cfg := &config.Config{BaseFile: base, OtherFile: other}
	var buf bytes.Buffer
	code := RunValidate(cfg, false, &buf)
	if code != 2 {
		t.Errorf("expected exit code 2 for changed entry with empty base, got %d", code)
	}
	if !stringContains(buf.String(), "SECRET") {
		t.Errorf("expected SECRET in output, got: %s", buf.String())
	}
}
