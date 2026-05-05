package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParse_BasicKeyValues(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDB_HOST=localhost\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(env.Entries))
	}
	if env.Index["APP_ENV"].Value != "production" {
		t.Errorf("APP_ENV: got %q, want %q", env.Index["APP_ENV"].Value, "production")
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret value"` + "\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.Index["SECRET"].Value != "my secret value" {
		t.Errorf("SECRET: got %q, want %q", env.Index["SECRET"].Value, "my secret value")
	}
}

func TestParse_InlineComment(t *testing.T) {
	path := writeTempEnv(t, "PORT=8080 # http port\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := env.Index["PORT"]
	if e.Value != "8080" {
		t.Errorf("PORT value: got %q, want %q", e.Value, "8080")
	}
	if e.Comment != "http port" {
		t.Errorf("PORT comment: got %q, want %q", e.Comment, "http port")
	}
}

func TestParse_SkipsCommentsAndBlanks(t *testing.T) {
	content := "# this is a comment\n\nKEY=value\n"
	path := writeTempEnv(t, content)
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(env.Entries))
	}
}

func TestParse_ExportPrefix(t *testing.T) {
	path := writeTempEnv(t, "export EXPORTED_VAR=hello\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := env.Index["EXPORTED_VAR"]; !ok {
		t.Error("expected EXPORTED_VAR to be parsed")
	}
}

func TestParse_MissingEquals(t *testing.T) {
	path := writeTempEnv(t, "BADLINE\n")
	_, err := Parse(path)
	if err == nil {
		t.Error("expected error for line missing '=', got nil")
	}
}

func TestParse_FileNotFound(t *testing.T) {
	_, err := Parse("/nonexistent/path/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
