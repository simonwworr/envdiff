package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/output"
	"github.com/user/envdiff/internal/parser"
)

func main() {
	maskSecrets := flag.Bool("mask", true, "mask sensitive values in output")
	summaryOnly := flag.Bool("summary", false, "print summary only")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [options] <file-a> <file-b>\n\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	fileA, fileB := args[0], args[1]

	envA, err := parseFile(fileA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", fileA, err)
		os.Exit(1)
	}

	envB, err := parseFile(fileB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", fileB, err)
		os.Exit(1)
	}

	results := diff.Compare(envA, envB)

	if *summaryOnly {
		output.Summary(os.Stdout, results)
	} else {
		output.WriteDiff(os.Stdout, results, *maskSecrets)
		fmt.Fprintln(os.Stdout)
		output.Summary(os.Stdout, results)
	}

	for _, r := range results {
		if r.Status != diff.Unchanged {
			os.Exit(1)
		}
	}
}

func parseFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parser.Parse(f)
}
