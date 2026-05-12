package reconcile

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// ExportFormat defines the output format for exported env vars.
type ExportFormat string

const (
	ExportFormatShell  ExportFormat = "shell"
	ExportFormatDocker ExportFormat = "docker"
	ExportFormatDotenv ExportFormat = "dotenv"
)

// ExportOptions controls how env vars are exported.
type ExportOptions struct {
	Format      ExportFormat
	OnlyChanged bool
	Masker      *diff.Masker
}

// Export writes env vars from a diff result to the given writer.
func Export(w io.Writer, entries []diff.Entry, opts ExportOptions) error {
	if opts.Format == "" {
		opts.Format = ExportFormatDotenv
	}

	filtered := filterEntries(entries, opts.OnlyChanged)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Key < filtered[j].Key
	})

	for _, e := range filtered {
		val := e.NewValue
		if val == "" {
			val = e.OldValue
		}
		if opts.Masker != nil && opts.Masker.IsSensitive(e.Key) {
			val = "***"
		}
		line, err := formatExportLine(e.Key, val, opts.Format)
		if err != nil {
			return err
		}
		fmt.Fprintln(w, line)
	}
	return nil
}

// ExportToFile writes the exported content to a file path.
func ExportToFile(path string, entries []diff.Entry, opts ExportOptions) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("export: cannot create file %q: %w", path, err)
	}
	defer f.Close()
	return Export(f, entries, opts)
}

func filterEntries(entries []diff.Entry, onlyChanged bool) []diff.Entry {
	if !onlyChanged {
		return entries
	}
	var out []diff.Entry
	for _, e := range entries {
		if e.Status != diff.StatusUnchanged {
			out = append(out, e)
		}
	}
	return out
}

func formatExportLine(key, val string, format ExportFormat) (string, error) {
	switch format {
	case ExportFormatShell:
		return fmt.Sprintf("export %s=%q", key, val), nil
	case ExportFormatDocker:
		return fmt.Sprintf("--env %s=%s", key, shellQuote(val)), nil
	case ExportFormatDotenv:
		return fmt.Sprintf("%s=%s", key, shellQuote(val)), nil
	default:
		return "", fmt.Errorf("export: unknown format %q", format)
	}
}

func shellQuote(s string) string {
	if strings.ContainsAny(s, " \t\n#$\"\'\\`") {
		return fmt.Sprintf("%q", s)
	}
	return s
}
