package reconcile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envdiff/internal/diff"
)

// Snapshot represents a point-in-time capture of a diff result.
type Snapshot struct {
	CreatedAt time.Time          `json:"created_at"`
	BaseFile  string             `json:"base_file"`
	OtherFile string             `json:"other_file"`
	Entries   []diff.Entry       `json:"entries"`
	Summary   SnapshotSummary    `json:"summary"`
}

// SnapshotSummary holds counts for a snapshot.
type SnapshotSummary struct {
	Added    int `json:"added"`
	Removed  int `json:"removed"`
	Changed  int `json:"changed"`
	Unchanged int `json:"unchanged"`
}

// TakeSnapshot creates a Snapshot from a diff result.
func TakeSnapshot(baseFile, otherFile string, result diff.Result) Snapshot {
	summary := SnapshotSummary{}
	for _, e := range result.Entries {
		switch e.Status {
		case diff.Added:
			summary.Added++
		case diff.Removed:
			summary.Removed++
		case diff.Changed:
			summary.Changed++
		case diff.Unchanged:
			summary.Unchanged++
		}
	}
	return Snapshot{
		CreatedAt: time.Now().UTC(),
		BaseFile:  baseFile,
		OtherFile: otherFile,
		Entries:   result.Entries,
		Summary:   summary,
	}
}

// SaveSnapshot writes a Snapshot as JSON to the given path.
func SaveSnapshot(path string, snap Snapshot) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create %q: %w", path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(snap)
}

// LoadSnapshot reads a Snapshot from a JSON file at the given path.
func LoadSnapshot(path string) (Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: open %q: %w", path, err)
	}
	defer f.Close()
	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: decode %q: %w", path, err)
	}
	return snap, nil
}
