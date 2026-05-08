package reconcile

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// SnapshotArgs holds CLI arguments for the snapshot subcommand.
type SnapshotArgs struct {
	BaseFile  string
	OtherFile string
	Output    string
}

// RunSnapshot parses two env files, diffs them, and writes a snapshot JSON.
// Returns a non-nil error on failure; writes a human-readable status to w.
func RunSnapshot(args SnapshotArgs, w io.Writer) error {
	base, err := parser.Parse(args.BaseFile)
	if err != nil {
		return fmt.Errorf("snapshot: parse base %q: %w", args.BaseFile, err)
	}
	other, err := parser.Parse(args.OtherFile)
	if err != nil {
		return fmt.Errorf("snapshot: parse other %q: %w", args.OtherFile, err)
	}

	result := diff.Compare(base, other)
	snap := TakeSnapshot(args.BaseFile, args.OtherFile, result)

	out := args.Output
	if out == "" {
		out = "envdiff-snapshot.json"
	}

	if err := SaveSnapshot(out, snap); err != nil {
		return err
	}

	fmt.Fprintf(w, "Snapshot saved to %s\n", out)
	fmt.Fprintf(w, "  added=%d removed=%d changed=%d unchanged=%d\n",
		snap.Summary.Added,
		snap.Summary.Removed,
		snap.Summary.Changed,
		snap.Summary.Unchanged,
	)
	return nil
}

// RunSnapshotFromArgs parses CLI flags and calls RunSnapshot.
func RunSnapshotFromArgs(arguments []string, w io.Writer) error {
	fs := flag.NewFlagSet("snapshot", flag.ContinueOnError)
	base := fs.String("base", "", "base .env file")
	other := fs.String("other", "", "other .env file")
	output := fs.String("output", "", "output snapshot path (default: envdiff-snapshot.json)")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(arguments); err != nil {
		return err
	}
	if *base == "" || *other == "" {
		return fmt.Errorf("snapshot: -base and -other are required")
	}
	return RunSnapshot(SnapshotArgs{
		BaseFile:  *base,
		OtherFile: *other,
		Output:    *output,
	}, w)
}
