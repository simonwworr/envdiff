package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeJSONResult(entries map[string]diff.Entry) diff.Result {
	return diff.Result{Entries: entries}
}

func TestJSONWriter_BasicOutput(t *testing.T) {
	result := makeJSONResult(map[string]diff.Entry{
		"FOO": {Status: diff.StatusAdded, NewValue: "bar"},
		"BAZ": {Status: diff.StatusRemoved, OldValue: "qux"},
	})

	var buf bytes.Buffer
	w := NewJSONWriter(&buf, nil, false)
	if err := w.Write(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	entries, ok := out["entries"].([]interface{})
	if !ok || len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %v", out["entries"])
	}
}

func TestJSONWriter_SummaryCorrect(t *testing.T) {
	result := makeJSONResult(map[string]diff.Entry{
		"A": {Status: diff.StatusAdded, NewValue: "1"},
		"B": {Status: diff.StatusRemoved, OldValue: "2"},
		"C": {Status: diff.StatusChanged, OldValue: "x", NewValue: "y"},
		"D": {Status: diff.StatusUnchanged, OldValue: "z", NewValue: "z"},
	})

	var buf bytes.Buffer
	w := NewJSONWriter(&buf, nil, true)
	if err := w.Write(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	type output struct {
		Summary struct {
			Added     int `json:"added"`
			Removed   int `json:"removed"`
			Changed   int `json:"changed"`
			Unchanged int `json:"unchanged"`
		} `json:"summary"`
	}
	var out output
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.Summary.Added != 1 || out.Summary.Removed != 1 || out.Summary.Changed != 1 || out.Summary.Unchanged != 1 {
		t.Errorf("unexpected summary: %+v", out.Summary)
	}
}

func TestJSONWriter_MasksSensitiveKeys(t *testing.T) {
	masker := diff.NewMasker()
	result := makeJSONResult(map[string]diff.Entry{
		"SECRET_KEY": {Status: diff.StatusChanged, OldValue: "oldpass", NewValue: "newpass"},
	})

	var buf bytes.Buffer
	w := NewJSONWriter(&buf, masker, true)
	if err := w.Write(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(buf.String(), "oldpass") || strings.Contains(buf.String(), "newpass") {
		t.Error("expected sensitive values to be masked in JSON output")
	}
}

func TestJSONWriter_IndentedOutput(t *testing.T) {
	result := makeJSONResult(map[string]diff.Entry{
		"X": {Status: diff.StatusUnchanged, OldValue: "1", NewValue: "1"},
	})

	var buf bytes.Buffer
	w := NewJSONWriter(&buf, nil, true)
	if err := w.Write(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "\n") {
		t.Error("expected indented (multi-line) JSON output")
	}
}
