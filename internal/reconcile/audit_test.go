package reconcile

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeAuditResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_NAME", Status: diff.StatusUnchanged, BaseValue: "myapp", OtherValue: "myapp"},
		{Key: "DB_HOST", Status: diff.StatusChanged, BaseValue: "localhost", OtherValue: "prod-db"},
		{Key: "NEW_KEY", Status: diff.StatusAdded, BaseValue: "", OtherValue: "newval"},
		{Key: "OLD_KEY", Status: diff.StatusRemoved, BaseValue: "gone", OtherValue: ""},
		{Key: "DB_PASSWORD", Status: diff.StatusChanged, BaseValue: "secret1", OtherValue: "secret2"},
	}
}

func TestNewAuditLog_EntryCount(t *testing.T) {
	masker := diff.NewMasker()
	results := makeAuditResults()
	log := NewAuditLog("base.env", "prod.env", results, masker)
	if len(log.Entries) != len(results) {
		t.Errorf("expected %d entries, got %d", len(results), len(log.Entries))
	}
}

func TestNewAuditLog_MasksSecretKey(t *testing.T) {
	masker := diff.NewMasker()
	results := makeAuditResults()
	log := NewAuditLog("base.env", "prod.env", results, masker)

	for _, e := range log.Entries {
		if e.Key == "DB_PASSWORD" {
			if !e.Masked {
				t.Error("expected DB_PASSWORD to be masked")
			}
			if e.OldValue != "***" || e.NewValue != "***" {
				t.Errorf("expected masked values, got old=%q new=%q", e.OldValue, e.NewValue)
			}
			return
		}
	}
	t.Error("DB_PASSWORD entry not found")
}

func TestNewAuditLog_PlaintextNonSecret(t *testing.T) {
	masker := diff.NewMasker()
	results := makeAuditResults()
	log := NewAuditLog("base.env", "prod.env", results, masker)

	for _, e := range log.Entries {
		if e.Key == "DB_HOST" {
			if e.Masked {
				t.Error("DB_HOST should not be masked")
			}
			if e.OldValue != "localhost" || e.NewValue != "prod-db" {
				t.Errorf("unexpected values: %q %q", e.OldValue, e.NewValue)
			}
			return
		}
	}
	t.Error("DB_HOST entry not found")
}

func TestAuditLog_Summary_ContainsFiles(t *testing.T) {
	masker := diff.NewMasker()
	log := NewAuditLog("base.env", "prod.env", makeAuditResults(), masker)
	summary := log.Summary()
	if !strings.Contains(summary, "base.env") || !strings.Contains(summary, "prod.env") {
		t.Errorf("summary missing file names: %s", summary)
	}
}

func TestAuditLog_Summary_ContainsMaskedLabel(t *testing.T) {
	masker := diff.NewMasker()
	log := NewAuditLog("base.env", "prod.env", makeAuditResults(), masker)
	summary := log.Summary()
	if !strings.Contains(summary, "[masked]") {
		t.Errorf("expected [masked] in summary, got: %s", summary)
	}
}

func TestAuditLog_EmptyEntries(t *testing.T) {
	masker := diff.NewMasker()
	log := NewAuditLog("a.env", "b.env", []diff.Result{}, masker)
	if len(log.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(log.Entries))
	}
	summary := log.Summary()
	if !strings.Contains(summary, "a.env") {
		t.Error("summary should still contain file names")
	}
}
