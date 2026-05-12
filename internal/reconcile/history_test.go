package reconcile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeHistoryEntry(action string, added, removed, changed int) HistoryEntry {
	return HistoryEntry{
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Action:    action,
		BaseFile:  ".env",
		OtherFile: ".env.staging",
		Added:     added,
		Removed:   removed,
		Changed:   changed,
	}
}

func TestHistory_Add_IncreasesCount(t *testing.T) {
	h := &History{}
	h.Add(makeHistoryEntry("promote", 2, 0, 1))
	h.Add(makeHistoryEntry("reconcile", 0, 1, 0))
	if len(h.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(h.Entries))
	}
}

func TestHistory_Add_SetsTimestampIfZero(t *testing.T) {
	h := &History{}
	e := HistoryEntry{Action: "promote"}
	h.Add(e)
	if h.Entries[0].Timestamp.IsZero() {
		t.Error("expected timestamp to be set automatically")
	}
}

func TestHistory_Latest_ReturnsLast(t *testing.T) {
	h := &History{}
	h.Add(makeHistoryEntry("reconcile", 1, 0, 0))
	h.Add(makeHistoryEntry("promote", 3, 1, 2))
	latest := h.Latest()
	if latest == nil || latest.Action != "promote" {
		t.Errorf("expected latest action 'promote', got %v", latest)
	}
}

func TestHistory_Latest_NilWhenEmpty(t *testing.T) {
	h := &History{}
	if h.Latest() != nil {
		t.Error("expected nil for empty history")
	}
}

func TestHistory_SortedByTime(t *testing.T) {
	h := &History{}
	e1 := makeHistoryEntry("a", 0, 0, 0)
	e1.Timestamp = time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)
	e2 := makeHistoryEntry("b", 0, 0, 0)
	e2.Timestamp = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	h.Add(e1)
	h.Add(e2)
	sorted := h.SortedByTime()
	if sorted[0].Action != "b" || sorted[1].Action != "a" {
		t.Errorf("unexpected sort order: %v %v", sorted[0].Action, sorted[1].Action)
	}
}

func TestSaveAndLoadHistory_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")
	h := &History{}
	h.Add(makeHistoryEntry("promote", 2, 1, 3))
	if err := SaveHistory(path, h); err != nil {
		t.Fatalf("SaveHistory: %v", err)
	}
	loaded, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Added != 2 {
		t.Errorf("expected Added=2, got %d", loaded.Entries[0].Added)
	}
}

func TestLoadHistory_MissingFile(t *testing.T) {
	h, err := LoadHistory("/nonexistent/path/history.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(h.Entries) != 0 {
		t.Errorf("expected empty history for missing file")
	}
}

func TestSaveHistory_ValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")
	h := &History{}
	h.Add(makeHistoryEntry("reconcile", 0, 0, 1))
	_ = SaveHistory(path, h)
	data, _ := os.ReadFile(path)
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Errorf("saved file is not valid JSON: %v", err)
	}
}
