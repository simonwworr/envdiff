package output

import (
	"io"

	"github.com/user/envdiff/internal/diff"
)

// Pipeline chains a masker and a Writer so callers only need to call Run.
type Pipeline struct {
	writer Writer
	masker *diff.Masker
	enabled bool
}

// NewPipeline constructs a Pipeline for the given format. When mask is true
// sensitive values are redacted before being handed to the writer.
func NewPipeline(format string, mask bool) *Pipeline {
	var m *diff.Masker
	if mask {
		m = diff.NewMasker()
	}
	return &Pipeline{
		writer:  NewWriter(format, mask),
		masker:  m,
		enabled: mask,
	}
}

// Run applies optional masking and writes the diff results to w.
func (p *Pipeline) Run(w io.Writer, results []diff.Result) error {
	out := results
	if p.enabled && p.masker != nil {
		out = make([]diff.Result, len(results))
		for i, r := range results {
			out[i] = p.masker.MaskResult(r)
		}
	}
	return p.writer.Write(w, out)
}
