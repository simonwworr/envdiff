package output

import (
	"io"

	"github.com/user/envdiff/internal/diff"
)

// Writer is the interface implemented by all output formatters.
type Writer interface {
	Write(w io.Writer, results []diff.Result) error
}

// Format constants mirror config.Format but avoid import cycles.
const (
	FormatText  = "text"
	FormatTable = "table"
	FormatJSON  = "json"
)

// NewWriter returns the appropriate Writer implementation for the given format
// string. An unknown format falls back to the plain-text writer.
func NewWriter(format string, masked bool) Writer {
	switch format {
	case FormatJSON:
		return NewJSONWriter(masked)
	case FormatTable:
		return NewTableWriter(masked)
	default:
		return &textWriter{masked: masked}
	}
}

// textWriter wraps the existing WriteDiff / Summary helpers as a Writer.
type textWriter struct {
	masked bool
}

func (t *textWriter) Write(w io.Writer, results []diff.Result) error {
	WriteDiff(w, results, t.masked)
	Summary(w, results)
	return nil
}
