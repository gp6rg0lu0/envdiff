package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// RenderEnvSet renders the keys of an EnvSet to w in the given format.
func RenderEnvSet(w io.Writer, s diff.EnvSet, format string) error {
	switch format {
	case "json":
		return renderEnvSetJSON(w, s)
	default:
		return renderEnvSetText(w, s)
	}
}

func renderEnvSetText(w io.Writer, s diff.EnvSet) error {
	keys := s.Keys()
	if len(keys) == 0 {
		_, err := fmt.Fprintln(w, "(empty set)")
		return err
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Keys (%d):\n", s.Len()))
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("  %s\n", k))
	}
	_, err := fmt.Fprint(w, sb.String())
	return err
}

func renderEnvSetJSON(w io.Writer, s diff.EnvSet) error {
	payload := struct {
		Count int      `json:"count"`
		Keys  []string `json:"keys"`
	}{
		Count: s.Len(),
		Keys:  s.Keys(),
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
