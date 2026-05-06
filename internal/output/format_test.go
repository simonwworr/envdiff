package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/output"
)

func makeResult(status diff.Status, key, aVal, bVal string) diff.Result {
	return diff.Result{
		Key:    key,
		Status: status,
		AValue: aVal,
		BValue: bVal,
	}
}

func TestWriteDiff_Added(t *testing.T) {
	var buf bytes.Buffer
	results := []diff.Result{makeResult(diff.Added, "NEW_KEY", "", "value")}
	output.WriteDiff(&buf, results, false)
	out := buf.String()
	if !strings.Contains(out, "NEW_KEY") {
		t.Errorf("expected NEW_KEY in output, got: %s", out)
	}
	if !strings.Contains(out, "+") {
		t.Errorf("expected '+' prefix for added key, got: %s", out)
	}
}

func TestWriteDiff_Removed(t *testing.T) {
	var buf bytes.Buffer
	results := []diff.Result{makeResult(diff.Removed, "OLD_KEY", "val", "")}
	output.WriteDiff(&buf, results, false)
	out := buf.String()
	if !strings.Contains(out, "-") {
		t.Errorf("expected '-' prefix for removed key, got: %s", out)
	}
}

func TestWriteDiff_Masked(t *testing.T) {
	var buf bytes.Buffer
	results := []diff.Result{makeResult(diff.Changed, "SECRET_TOKEN", "abc", "xyz")}
	output.WriteDiff(&buf, results, true)
	out := buf.String()
	if strings.Contains(out, "abc") || strings.Contains(out, "xyz") {
		t.Errorf("expected secret values to be masked, got: %s", out)
	}
}

func TestSummary_Counts(t *testing.T) {
	var buf bytes.Buffer
	results := []diff.Result{
		makeResult(diff.Added, "A", "", "1"),
		makeResult(diff.Removed, "B", "2", ""),
		makeResult(diff.Changed, "C", "3", "4"),
		makeResult(diff.Unchanged, "D", "5", "5"),
	}
	output.Summary(&buf, results)
	out := buf.String()
	for _, want := range []string{"1 added", "1 removed", "1 changed"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in summary output, got: %s", want, out)
		}
	}
}
