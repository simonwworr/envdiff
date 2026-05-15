package reconcile

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeRedactResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_NAME", BaseValue: "myapp", OtherValue: "myapp", Status: diff.Unchanged},
		{Key: "DB_PASSWORD", BaseValue: "s3cr3t", OtherValue: "s3cr3t", Status: diff.Unchanged},
		{Key: "API_KEY", BaseValue: "", OtherValue: "abc123", Status: diff.Added},
		{Key: "PORT", BaseValue: "8080", OtherValue: "9090", Status: diff.Changed},
	}
}

func TestRedact_TotalKeyCount(t *testing.T) {
	results := makeRedactResults()
	masker := diff.NewMasker()
	report := Redact(results, masker)

	if report.TotalKeys != 4 {
		t.Errorf("expected 4 total keys, got %d", report.TotalKeys)
	}
}

func TestRedact_MasksSensitiveKeys(t *testing.T) {
	results := makeRedactResults()
	masker := diff.NewMasker()
	report := Redact(results, masker)

	for _, entry := range report.Entries {
		if masker.IsSensitive(entry.Key) && entry.Original != "" {
			if entry.Redacted != "[REDACTED]" {
				t.Errorf("key %q should be redacted, got %q", entry.Key, entry.Redacted)
			}
			if !entry.WasMasked {
				t.Errorf("key %q WasMasked should be true", entry.Key)
			}
		}
	}
}

func TestRedact_PreservesNonSensitiveKeys(t *testing.T) {
	results := makeRedactResults()
	masker := diff.NewMasker()
	report := Redact(results, masker)

	for _, entry := range report.Entries {
		if !masker.IsSensitive(entry.Key) {
			if entry.WasMasked {
				t.Errorf("key %q should not be masked", entry.Key)
			}
			if entry.Redacted != entry.Original {
				t.Errorf("key %q value should be unchanged", entry.Key)
			}
		}
	}
}

func TestRedact_RedactedKeyCount(t *testing.T) {
	results := makeRedactResults()
	masker := diff.NewMasker()
	report := Redact(results, masker)

	if report.RedactedKeys == 0 {
		t.Error("expected at least one redacted key")
	}
}

func TestRedactReport_Summary(t *testing.T) {
	results := makeRedactResults()
	masker := diff.NewMasker()
	report := Redact(results, masker)

	summary := report.Summary()
	if !strings.Contains(summary, "redacted") {
		t.Errorf("summary should contain 'redacted', got %q", summary)
	}
}

func TestRedact_NilMasker(t *testing.T) {
	results := makeRedactResults()
	report := Redact(results, nil)

	if report.RedactedKeys != 0 {
		t.Errorf("expected 0 redacted keys with nil masker, got %d", report.RedactedKeys)
	}
}
