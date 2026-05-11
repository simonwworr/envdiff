package reconcile

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// AuditArgs holds parsed CLI arguments for the audit subcommand.
type AuditArgs struct {
	BaseFile  string
	OtherFile string
	JSONOut   bool
	OutFile   string
}

// ParseAuditArgs parses flags for the audit subcommand.
func ParseAuditArgs(args []string) (AuditArgs, error) {
	fs := flag.NewFlagSet("audit", flag.ContinueOnError)
	jsonOut := fs.Bool("json", false, "output audit log as JSON")
	outFile := fs.String("out", "", "write audit log to file instead of stdout")
	if err := fs.Parse(args); err != nil {
		return AuditArgs{}, err
	}
	if fs.NArg() < 2 {
		return AuditArgs{}, fmt.Errorf("usage: audit [flags] <base.env> <other.env>")
	}
	return AuditArgs{
		BaseFile:  fs.Arg(0),
		OtherFile: fs.Arg(1),
		JSONOut:   *jsonOut,
		OutFile:   *outFile,
	}, nil
}

// RunAudit executes the audit command writing output to w.
func RunAudit(args AuditArgs, w io.Writer) error {
	baseEnv, err := parser.Parse(args.BaseFile)
	if err != nil {
		return fmt.Errorf("parse base: %w", err)
	}
	otherEnv, err := parser.Parse(args.OtherFile)
	if err != nil {
		return fmt.Errorf("parse other: %w", err)
	}

	results := diff.Compare(baseEnv, otherEnv)
	masker := diff.NewMasker()
	log := NewAuditLog(args.BaseFile, args.OtherFile, results, masker)

	out := w
	if args.OutFile != "" {
		f, err := os.Create(args.OutFile)
		if err != nil {
			return fmt.Errorf("open output file: %w", err)
		}
		defer f.Close()
		out = f
	}

	if args.JSONOut {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(log)
	}
	_, err = fmt.Fprint(out, log.Summary())
	return err
}
