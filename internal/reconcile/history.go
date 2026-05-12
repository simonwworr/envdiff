package reconcile

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// HistoryEntry records a single reconcile or promote action applied to an env file.
type HistoryEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	BaseFile  string    `json:"base_file"`
	OtherFile string    `json:"other_file,omitempty"`
	Added     int       `json:"added"`
	Removed   int       `json:"removed"`
	Changed   int       `json:"changed"`
	Note      string    `json:"note,omitempty"`
}

// History holds an ordered list of history entries.
type History struct {
	Entries []HistoryEntry `json:"entries"`
}

// Add appends a new entry to the history.
func (h *History) Add(entry HistoryEntry) {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	h.Entries = append(h.Entries, entry)
}

// Latest returns the most recent entry, or nil if history is empty.
func (h *History) Latest() *HistoryEntry {
	if len(h.Entries) == 0 {
		return nil
	}
	return &h.Entries[len(h.Entries)-1]
}

// SortedByTime returns entries sorted ascending by timestamp.
func (h *History) SortedByTime() []HistoryEntry {
	copy := append([]HistoryEntry{}, h.Entries...)
	sort.Slice(copy, func(i, j int) bool {
		return copy[i].Timestamp.Before(copy[j].Timestamp)
	})
	return copy
}

// SaveHistory writes the history to a JSON file at path.
func SaveHistory(path string, h *History) error {
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal history: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadHistory reads a history file from path.
// Returns an empty History if the file does not exist.
func LoadHistory(path string) (*History, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &History{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read history: %w", err)
	}
	var h History
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, fmt.Errorf("unmarshal history: %w", err)
	}
	return &h, nil
}
