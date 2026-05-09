package reconcile

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempPromoteEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestRunPromote_BasicPromotion(t *testing.T) {
	src := writeTempPromoteEnv(t, "FOO=bar\nBAZ=qux\n")
	dst := writeTempPromoteEnv(t, "EXISTING=yes\n")

	pa := PromoteArgs{SourceFile: src, TargetFile: dst}
	var buf bytes.Buffer
	if err := RunPromote(pa, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "applied=2") {
		t.Errorf("expected applied=2 in output, got: %q", out)
	}
}

func TestRunPromote_SkipSensitiveKeys(t *testing.T) {
	src := writeTempPromoteEnv(t, "API_SECRET=topsecret\nAPP_NAME=myapp\n")
	dst := writeTempPromoteEnv(t, "")

	pa := PromoteArgs{SourceFile: src, TargetFile: dst, SkipSensitive: true}
	var buf bytes.Buffer
	if err := RunPromote(pa, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "skipped=1") {
		t.Errorf("expected skipped=1 in output, got: %q", out)
	}
}

func TestRunPromote_WritesToFile(t *testing.T) {
	src := writeTempPromoteEnv(t, "FOO=bar\n")
	dst := writeTempPromoteEnv(t, "")
	outFile := filepath.Join(t.TempDir(), "promoted.env")

	pa := PromoteArgs{SourceFile: src, TargetFile: dst, OutputFile: outFile}
	var buf bytes.Buffer
	if err := RunPromote(pa, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("output file not written: %v", err)
	}
	if !strings.Contains(string(data), "FOO") {
		t.Errorf("expected FOO in output file, got: %q", string(data))
	}
}

func TestRunPromote_MissingSourceFile(t *testing.T) {
	pa := PromoteArgs{SourceFile: "/nonexistent.env", TargetFile: "/also/missing.env"}
	var buf bytes.Buffer
	err := RunPromote(pa, &buf)
	if err == nil {
		t.Fatal("expected error for missing source file")
	}
}

func TestParsePromoteArgs_MissingPositional(t *testing.T) {
	_, err := ParsePromoteArgs([]string{"--overwrite"})
	if err == nil {
		t.Fatal("expected error when positional args missing")
	}
}
