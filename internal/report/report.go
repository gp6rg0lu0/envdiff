package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Options controls report rendering behavior.
type Options struct {
	Format  Format
	Redact  bool
	NoColor bool
}

// Render writes a human-readable or JSON report of diff results to w.
func Render(w io.Writer, result diff.Result, baseFile, compareFile string, opts Options) error {
	if opts.Redact {
		result = diff.RedactResult(result)
	}

	switch opts.Format {
	case FormatJSON:
		return renderJSON(w, result, baseFile, compareFile)
	default:
		return renderText(w, result, baseFile, compareFile, opts.NoColor)
	}
}

func renderText(w io.Writer, result diff.Result, baseFile, compareFile string, noColor bool) error {
	red := colorFn("\033[31m", noColor)
	green := colorFn("\033[32m", noColor)
	yellow := colorFn("\033[33m", noColor)
	reset := colorFn("\033[0m", noColor)

	if len(result.Missing) == 0 && len(result.Extra) == 0 && len(result.Mismatched) == 0 {
		fmt.Fprintln(w, green("✓ No differences found."+reset("")))
		return nil
	}

	fmt.Fprintf(w, "Comparing %s → %s\n", baseFile, compareFile)
	fmt.Fprintln(w, strings.Repeat("-", 40))

	for _, k := range result.Missing {
		fmt.Fprintf(w, "%s MISSING  %s%s\n", red("✗"), k, reset(""))
	}
	for _, k := range result.Extra {
		fmt.Fprintf(w, "%s EXTRA    %s%s\n", green("+"), k, reset(""))
	}
	for _, m := range result.Mismatched {
		fmt.Fprintf(w, "%s MISMATCH %s: %q → %q%s\n", yellow("~"), m.Key, m.BaseValue, m.CompareValue, reset(""))
	}

	return nil
}

func colorFn(code string, noColor bool) func(string) string {
	if noColor {
		return func(s string) string { return s }
	}
	return func(s string) string { return code + s }
}
