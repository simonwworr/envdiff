package diff

import (
	"fmt"
	"sort"
)

// EntryStatus represents the diff status of a key.
type EntryStatus string

const (
	StatusAdded    EntryStatus = "added"
	StatusRemoved  EntryStatus = "removed"
	StatusChanged  EntryStatus = "changed"
	StatusUnchanged EntryStatus = "unchanged"
)

// Entry represents a single diff result for a key.
type Entry struct {
	Key      string
	Status   EntryStatus
	BaseVal  string
	OtherVal string
}

// Result holds the full diff between two env maps.
type Result struct {
	Entries []Entry
}

// HasChanges returns true if there are any non-unchanged entries.
func (r *Result) HasChanges() bool {
	for _, e := range r.Entries {
		if e.Status != StatusUnchanged {
			return true
		}
	}
	return false
}

// Compare diffs base against other, returning a Result.
// Keys present only in base are Removed; only in other are Added;
// in both with different values are Changed; same are Unchanged.
func Compare(base, other map[string]string) *Result {
	keys := unionKeys(base, other)
	sort.Strings(keys)

	var entries []Entry
	for _, k := range keys {
		bv, inBase := base[k]
		ov, inOther := other[k]

		var status EntryStatus
		switch {
		case inBase && !inOther:
			status = StatusRemoved
		case !inBase && inOther:
			status = StatusAdded
		case bv != ov:
			status = StatusChanged
		default:
			status = StatusUnchanged
		}

		entries = append(entries, Entry{
			Key:      k,
			Status:   status,
			BaseVal:  bv,
			OtherVal: ov,
		})
	}
	return &Result{Entries: entries}
}

// Summary returns a human-readable one-line summary of the diff.
func (r *Result) Summary() string {
	var added, removed, changed int
	for _, e := range r.Entries {
		switch e.Status {
		case StatusAdded:
			added++
		case StatusRemoved:
			removed++
		case StatusChanged:
			changed++
		}
	}
	return fmt.Sprintf("+%d added, -%d removed, ~%d changed", added, removed, changed)
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
