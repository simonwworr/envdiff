package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makePipelineResults() []diff.Result {
	return []diff.Result{
		{Key: "API_SECRET", Status: diff.Changed, BaseValue: "old-secret", OtherValue: "new-secret"},
		{Key: "HOST", Status: diff.Unchanged, BaseValue: "localhost", OtherValue: "localhost"},
	}
}

func TestNewPipeline_NoMask(t *testing.T) {
	p := NewPipeline(FormatText, false)
	if p.enabled {
		t.Error("expected masking disabled")
	}
	if p.masker != nil {
		t.Error("expected nil masker when mask=false")
	}
}

func TestNewPipeline_WithMask(t *testing.T) {
	p := NewPipeline(FormatText, true)
	if !p.enabled {
		t.Error("expected masking enabled")
	}
	if p.masker == nil {
		t.Error("expected non-nil masker when mask=true")
	}
}

func TestPipeline_Run_MasksSecret(t *testing.T) {
	var buf bytes.Buffer
	p := NewPipeline(FormatText, true)
	if err := p.Run(&buf, makePipelineResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "old-secret") || strings.Contains(out, "new-secret") {
		t.Error("expected secret values to be masked in output")
	}
}

func TestPipeline_Run_NoMask_ShowsValues(t *testing.T) {
	var buf bytes.Buffer
	p := NewPipeline(FormatText, false)
	if err := p.Run(&buf, makePipelineResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "old-secret") {
		t.Error("expected unmasked secret value in output when mask=false")
	}
}

func TestPipeline_Run_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	p := NewPipeline(FormatJSON, false)
	if err := p.Run(&buf, makePipelineResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(strings.TrimSpace(buf.String()), "{") {
		t.Error("expected JSON output to start with '{'")
	}
}
