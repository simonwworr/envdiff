package diff

import (
	"testing"
)

func TestCompare_Added(t *testing.T) {
	base := map[string]string{"A": "1"}
	other := map[string]string{"A": "1", "B": "2"}

	r := Compare(base, other)
	entry := findEntry(r, "B")
	if entry == nil {
		t.Fatal("expected entry for B")
	}
	if entry.Status != StatusAdded {
		t.Errorf("expected added, got %s", entry.Status)
	}
	if entry.OtherVal != "2" {
		t.Errorf("expected OtherVal=2, got %s", entry.OtherVal)
	}
}

func TestCompare_Removed(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	other := map[string]string{"A": "1"}

	r := Compare(base, other)
	entry := findEntry(r, "B")
	if entry == nil {
		t.Fatal("expected entry for B")
	}
	if entry.Status != StatusRemoved {
		t.Errorf("expected removed, got %s", entry.Status)
	}
}

func TestCompare_Changed(t *testing.T) {
	base := map[string]string{"A": "old"}
	other := map[string]string{"A": "new"}

	r := Compare(base, other)
	entry := findEntry(r, "A")
	if entry == nil {
		t.Fatal("expected entry for A")
	}
	if entry.Status != StatusChanged {
		t.Errorf("expected changed, got %s", entry.Status)
	}
	if entry.BaseVal != "old" || entry.OtherVal != "new" {
		t.Errorf("unexpected values: base=%s other=%s", entry.BaseVal, entry.OtherVal)
	}
}

func TestCompare_Unchanged(t *testing.T) {
	base := map[string]string{"A": "same"}
	other := map[string]string{"A": "same"}

	r := Compare(base, other)
	entry := findEntry(r, "A")
	if entry == nil {
		t.Fatal("expected entry for A")
	}
	if entry.Status != StatusUnchanged {
		t.Errorf("expected unchanged, got %s", entry.Status)
	}
}

func TestResult_HasChanges(t *testing.T) {
	r := &Result{Entries: []Entry{{Key: "X", Status: StatusAdded}}}
	if !r.HasChanges() {
		t.Error("expected HasChanges=true")
	}

	r2 := &Result{Entries: []Entry{{Key: "X", Status: StatusUnchanged}}}
	if r2.HasChanges() {
		t.Error("expected HasChanges=false")
	}
}

func TestResult_Summary(t *testing.T) {
	r := Compare(
		map[string]string{"A": "1", "B": "old"},
		map[string]string{"B": "new", "C": "3"},
	)
	summary := r.Summary()
	expected := "+1 added, -1 removed, ~1 changed"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}

func findEntry(r *Result, key string) *Entry {
	for i := range r.Entries {
		if r.Entries[i].Key == key {
			return &r.Entries[i]
		}
	}
	return nil
}
