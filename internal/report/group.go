package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
)

// RenderGrouped writes a grouped view of the diff result to w.
func RenderGrouped(w io.Writer, r diff.Result, format string) error {
	g := diff.Group(r)
	switch format {
	case "json":
		return renderGroupedJSON(w, g)
	default:
		return renderGroupedText(w, g)
	}
}

func renderGroupedText(w io.Writer, g diff.GroupedResult) error {
	if g.Total() == 0 {
		_, err := fmt.Fprintln(w, "No differences found.")
		return err
	}

	if len(g.Missing) > 0 {
		fmt.Fprintf(w, "Missing keys (%d):\n", len(g.Missing))
		for _, e := range g.Missing {
			fmt.Fprintf(w, "  - %s\n", e.Key)
		}
	}

	if len(g.Extra) > 0 {
		fmt.Fprintf(w, "Extra keys (%d):\n", len(g.Extra))
		for _, e := range g.Extra {
			fmt.Fprintf(w, "  + %s\n", e.Key)
		}
	}

	if len(g.Mismatched) > 0 {
		fmt.Fprintf(w, "Mismatched keys (%d):\n", len(g.Mismatched))
		for _, e := range g.Mismatched {
			fmt.Fprintf(w, "  ~ %s: %q vs %q\n", e.Key, e.BaseValue, e.OtherValue)
		}
	}

	fmt.Fprintf(w, "Total: %d difference(s)\n", g.Total())
	return nil
}

func renderGroupedJSON(w io.Writer, g diff.GroupedResult) error {
	type payload struct {
		Missing    []string            `json:"missing"`
		Extra      []string            `json:"extra"`
		Mismatched []map[string]string `json:"mismatched"`
		Total      int                 `json:"total"`
	}

	p := payload{Total: g.Total()}
	for _, e := range g.Missing {
		p.Missing = append(p.Missing, e.Key)
	}
	for _, e := range g.Extra {
		p.Extra = append(p.Extra, e.Key)
	}
	for _, e := range g.Mismatched {
		p.Mismatched = append(p.Mismatched, map[string]string{
			"key":   e.Key,
			"base":  e.BaseValue,
			"other": e.OtherValue,
		})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}
