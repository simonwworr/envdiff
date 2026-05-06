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
	baseFile := flag.String("base", "", "Base .env file (required)")
	targetFile := flag.String("target", "", "Target .env file (required)")
	maskSecrets := flag.Bool("mask", true, "Mask sensitive values")
	tableFormat := flag.Bool("table", false, "Render output as aligned table")
	flag.Parse()

	if *baseFile == "" || *targetFile == "" {
		fmt.Fprintln(os.Stderr, "Usage: envdiff -base <file> -target <file> [-mask] [-table]")
		os.Exit(1)
	}

	base, err := parseFile(*baseFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading base file: %v\n", err)
		os.Exit(1)
	}

	target, err := parseFile(*targetFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading target file: %v\n", err)
		os.Exit(1)
	}

	result := diff.Compare(base, target)

	var masker *diff.Masker
	if *maskSecrets {
		masker = diff.NewMasker()
	}

	if *tableFormat {
		tw := output.NewTableWriter(os.Stdout, masker)
		if err := tw.Write(result); err != nil {
			fmt.Fprintf(os.Stderr, "error writing table: %v\n", err)
			os.Exit(1)
		}
	} else {
		output.WriteDiff(os.Stdout, result, masker)
	}

	fmt.Fprintln(os.Stdout)
	output.Summary(os.Stdout, result)
}

func parseFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parser.Parse(f)
}
