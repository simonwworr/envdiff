package reconcile

import (
	"strings"
	"testing"
)

func TestMerge_OnlyInBase(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	other := map[string]string{}
	r := Merge(base, other, PreferBase)
	if r.Env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", r.Env["FOO"])
	}
	if len(r.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(r.Conflicts))
	}
}

func TestMerge_OnlyInOther(t *testing.T) {
	base := map[string]string{}
	other := map[string]string{"BAZ": "qux"}
	r := Merge(base, other, PreferBase)
	if r.Env["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", r.Env["BAZ"])
	}
}

func TestMerge_NoConflict(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	other := map[string]string{"A": "1", "C": "3"}
	r := Merge(base, other, PreferBase)
	if r.Env["A"] != "1" || r.Env["B"] != "2" || r.Env["C"] != "3" {
		t.Errorf("unexpected env: %v", r.Env)
	}
	if len(r.Conflicts) != 0 {
		t.Errorf("expected no conflicts")
	}
}

func TestMerge_ConflictPreferBase(t *testing.T) {
	base := map[string]string{"KEY": "base_val"}
	other := map[string]string{"KEY": "other_val"}
	r := Merge(base, other, PreferBase)
	if r.Env["KEY"] != "base_val" {
		t.Errorf("expected base_val, got %q", r.Env["KEY"])
	}
	if len(r.Conflicts) != 1 || r.Conflicts[0].Chosen != "base_val" {
		t.Errorf("unexpected conflicts: %+v", r.Conflicts)
	}
}

func TestMerge_ConflictPreferOther(t *testing.T) {
	base := map[string]string{"KEY": "base_val"}
	other := map[string]string{"KEY": "other_val"}
	r := Merge(base, other, PreferOther)
	if r.Env["KEY"] != "other_val" {
		t.Errorf("expected other_val, got %q", r.Env["KEY"])
	}
	if len(r.Conflicts) != 1 || r.Conflicts[0].Chosen != "other_val" {
		t.Errorf("unexpected conflicts: %+v", r.Conflicts)
	}
}

func TestConflictSummary_NoConflicts(t *testing.T) {
	s := ConflictSummary(nil)
	if s != "no conflicts" {
		t.Errorf("expected 'no conflicts', got %q", s)
	}
}

func TestConflictSummary_WithConflicts(t *testing.T) {
	conflicts := []MergeConflict{
		{Key: "DB_PASS", BaseVal: "secret", OtherVal: "other", Chosen: "secret"},
	}
	s := ConflictSummary(conflicts)
	if !strings.Contains(s, "1 conflict") {
		t.Errorf("expected conflict count in summary, got %q", s)
	}
	if !strings.Contains(s, "DB_PASS") {
		t.Errorf("expected key name in summary, got %q", s)
	}
}
