package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
)

// RenderBaseline renders the diff between a baseline snapshot and the current
// environment in the requested format.
func RenderBaseline(w io.Writer, baseline diff.Baseline, result diff.Result, format string, redact bool) error {
	if redact {
		result = diff.RedactResult(result)
	}
	switch format {
	case "json":
		return renderBaselineJSON(w, baseline, result)
	default:
		return renderBaselineText(w, baseline, result)
	}
}

func renderBaselineText(w io.Writer, baseline diff.Baseline, result diff.Result) error {
	fmt.Fprintf(w, "Baseline: %s (%d keys)\n", baseline.Name, len(baseline.Entries))
	if len(result.Missing) == 0 && len(result.Extra) == 0 && len(result.Mismatched) == 0 {
		fmt.Fprintln(w, "No changes since baseline.")
		return nil
	}
	for _, k := range result.Missing {
		fmt.Fprintf(w, "  REMOVED  %s\n", k)
	}
	for _, k := range result.Extra {
		fmt.Fprintf(w, "  ADDED    %s\n", k)
	}
	for k, v := range result.Mismatched {
		fmt.Fprintf(w, "  CHANGED  %s: %q -> %q\n", k, v.Base, v.Compare)
	}
	return nil
}

func renderBaselineJSON(w io.Writer, baseline diff.Baseline, result diff.Result) error {
	type payload struct {
		Baseline string      `json:"baseline"`
		Keys     int         `json:"baseline_keys"`
		Result   diff.Result `json:"changes"`
	}
	p := payload{
		Baseline: baseline.Name,
		Keys:     len(baseline.Entries),
		Result:   result,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}
