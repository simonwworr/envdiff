package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeTableResult(entries []diff.Entry) diff.Result {
	return diff.Result{Entries: entries}
}

func TestTableWriter_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	tw := NewTableWriter(&buf, nil)
	err := tw.Write(makeTableResult(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestTableWriter_HeaderPresent(t *testing.T) {
	var buf bytes.Buffer
	tw := NewTableWriter(&buf, nil)
	entries := []diff.Entry{
		{Key: "APP_ENV", Status: diff.StatusAdded, BaseValue: "", TargetValue: "production"},
	}
	tw.Write(makeTableResult(entries))
	out := buf.String()
	for _, col := range []string{"KEY", "STATUS", "BASE", "TARGET"} {
		if !strings.Contains(out, col) {
			t.Errorf("expected header column %q in output", col)
		}
	}
}

func TestTableWriter_RowContent(t *testing.T) {
	var buf bytes.Buffer
	tw := NewTableWriter(&buf, nil)
	entries := []diff.Entry{
		{Key: "DB_HOST", Status: diff.StatusChanged, BaseValue: "localhost", TargetValue: "db.prod"},
	}
	tw.Write(makeTableResult(entries))
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected key in output")
	}
	if !strings.Contains(out, "localhost") {
		t.Errorf("expected base value in output")
	}
	if !strings.Contains(out, "db.prod") {
		t.Errorf("expected target value in output")
	}
}

func TestTableWriter_TruncatesLongValues(t *testing.T) {
	var buf bytes.Buffer
	tw := NewTableWriter(&buf, nil)
	long := strings.Repeat("x", 30)
	entries := []diff.Entry{
		{Key: "SECRET", Status: diff.StatusAdded, BaseValue: "", TargetValue: long},
	}
	tw.Write(makeTableResult(entries))
	out := buf.String()
	if strings.Contains(out, long) {
		t.Errorf("expected long value to be truncated")
	}
	if !strings.Contains(out, "...") {
		t.Errorf("expected ellipsis for truncated value")
	}
}

func TestTableWriter_MaskedSensitiveKey(t *testing.T) {
	var buf bytes.Buffer
	masker := diff.NewMasker()
	tw := NewTableWriter(&buf, masker)
	entries := []diff.Entry{
		{Key: "DB_PASSWORD", Status: diff.StatusChanged, BaseValue: "secret1", TargetValue: "secret2"},
	}
	tw.Write(makeTableResult(entries))
	out := buf.String()
	if strings.Contains(out, "secret1") || strings.Contains(out, "secret2") {
		t.Errorf("expected sensitive values to be masked")
	}
}
