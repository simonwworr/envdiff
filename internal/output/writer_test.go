package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeWriterResults() []diff.Result {
	return []diff.Result{
		{Key: "HOST", Status: diff.Added, OtherValue: "localhost"},
		{Key: "PORT", Status: diff.Removed, BaseValue: "8080"},
		{Key: "DEBUG", Status: diff.Changed, BaseValue: "false", OtherValue: "true"},
	}
}

func TestNewWriter_DefaultIsText(t *testing.T) {
	w := NewWriter("", false)
	if _, ok := w.(*textWriter); !ok {
		t.Errorf("expected *textWriter for empty format, got %T", w)
	}
}

func TestNewWriter_TextFormat(t *testing.T) {
	w := NewWriter(FormatText, false)
	if _, ok := w.(*textWriter); !ok {
		t.Errorf("expected *textWriter, got %T", w)
	}
}

func TestNewWriter_JSONFormat(t *testing.T) {
	w := NewWriter(FormatJSON, false)
	if _, ok := w.(*JSONWriter); !ok {
		t.Errorf("expected *JSONWriter, got %T", w)
	}
}

func TestNewWriter_TableFormat(t *testing.T) {
	w := NewWriter(FormatTable, false)
	if _, ok := w.(*TableWriter); !ok {
		t.Errorf("expected *TableWriter, got %T", w)
	}
}

func TestTextWriter_Write_ContainsKeys(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(FormatText, false)
	if err := w.Write(&buf, makeWriterResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, key := range []string{"HOST", "PORT", "DEBUG"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected output to contain key %q", key)
		}
	}
}

func TestTextWriter_Write_ContainsSummary(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(FormatText, false)
	_ = w.Write(&buf, makeWriterResults())
	if !strings.Contains(buf.String(), "added") {
		t.Error("expected summary line in text output")
	}
}
