package output

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// JSONWriter writes diff results as structured JSON output.
type JSONWriter struct {
	w       io.Writer
	masker  *diff.Masker
	indent  bool
}

// NewJSONWriter creates a JSONWriter that writes to w.
// If indent is true, output is pretty-printed.
func NewJSONWriter(w io.Writer, masker *diff.Masker, indent bool) *JSONWriter {
	return &JSONWriter{w: w, masker: masker, indent: indent}
}

type jsonEntry struct {
	Key      string `json:"key"`
	Status   string `json:"status"`
	OldValue string `json:"old_value,omitempty"`
	NewValue string `json:"new_value,omitempty"`
}

type jsonOutput struct {
	Entries []jsonEntry `json:"entries"`
	Summary jsonSummary `json:"summary"`
}

type jsonSummary struct {
	Added    int `json:"added"`
	Removed  int `json:"removed"`
	Changed  int `json:"changed"`
	Unchanged int `json:"unchanged"`
}

// Write serializes the diff result to JSON and writes it to the underlying writer.
func (j *JSONWriter) Write(result diff.Result) error {
	keys := make([]string, 0, len(result.Entries))
	for k := range result.Entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := jsonOutput{Entries: make([]jsonEntry, 0, len(keys))}

	for _, k := range keys {
		e := result.Entries[k]
		entry := jsonEntry{
			Key:    k,
			Status: string(e.Status),
		}
		if j.masker != nil && j.masker.IsSensitive(k) {
			entry.OldValue = maskIfNonEmpty(e.OldValue)
			entry.NewValue = maskIfNonEmpty(e.NewValue)
		} else {
			entry.OldValue = e.OldValue
			entry.NewValue = e.NewValue
		}
		out.Entries = append(out.Entries, entry)

		switch e.Status {
		case diff.StatusAdded:
			out.Summary.Added++
		case diff.StatusRemoved:
			out.Summary.Removed++
		case diff.StatusChanged:
			out.Summary.Changed++
		case diff.StatusUnchanged:
			out.Summary.Unchanged++
		}
	}

	var data []byte
	var err error
	if j.indent {
		data, err = json.MarshalIndent(out, "", "  ")
	} else {
		data, err = json.Marshal(out)
	}
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}
	_, err = fmt.Fprintln(j.w, string(data))
	return err
}
