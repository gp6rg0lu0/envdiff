package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
)

// RenderScore writes the compatibility score to w in the given format.
func RenderScore(w io.Writer, s diff.Score, format string) error {
	switch format {
	case "json":
		return renderScoreJSON(w, s)
	default:
		return renderScoreText(w, s)
	}
}

func renderScoreText(w io.Writer, s diff.Score) error {
	_, err := fmt.Fprintf(w, "Compatibility Score: %s\n", s.String())
	return err
}

type scoreJSON struct {
	Percent    float64 `json:"percent"`
	Total      int     `json:"total_keys"`
	Matched    int     `json:"matched"`
	Missing    int     `json:"missing"`
	Extra      int     `json:"extra"`
	Mismatched int     `json:"mismatched"`
}

func renderScoreJSON(w io.Writer, s diff.Score) error {
	payload := scoreJSON{
		Percent:    s.Percent,
		Total:      s.Total,
		Matched:    s.Matched,
		Missing:    s.Missing,
		Extra:      s.Extra,
		Mismatched: s.Mismatched,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
