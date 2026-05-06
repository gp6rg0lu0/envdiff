package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envdiff/internal/diff"
	"github.com/yourorg/envdiff/internal/parser"
	"github.com/yourorg/envdiff/internal/report"
)

func main() {
	var (
		redact bool
		format string
	)

	flag.BoolVar(&redact, "redact", false, "Redact sensitive values in output")
	flag.StringVar(&format, "format", "text", "Output format: text or json")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [options] <base.env> <compare.env>\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	baseFile := args[0]
	cmpFile := args[1]

	baseEnv, err := parser.ParseFile(baseFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading base file %q: %v\n", baseFile, err)
		os.Exit(1)
	}

	cmpEnv, err := parser.ParseFile(cmpFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading compare file %q: %v\n", cmpFile, err)
		os.Exit(1)
	}

	result := diff.Compare(baseEnv, cmpEnv)

	if redact {
		result = diff.RedactResult(result)
	}

	if err := report.Render(os.Stdout, result, format); err != nil {
		fmt.Fprintf(os.Stderr, "error rendering report: %v\n", err)
		os.Exit(1)
	}

	if len(result.Missing) > 0 || len(result.Extra) > 0 || len(result.Mismatched) > 0 {
		os.Exit(2)
	}
}
