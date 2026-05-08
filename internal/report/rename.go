package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// RenderRename renders a RenameResult as either text or JSON.
func RenderRename(result diff.RenameResult, format string) string {
	switch format {
	case "json":
		return renderRenameJSON(result)
	default:
		return renderRenameText(result)
	}
}

func renderRenameText(result diff.RenameResult) string {
	var sb strings.Builder

	if len(result.Renamed) == 0 && len(result.NotFound) == 0 {
		sb.WriteString("No renames applied.\n")
		return sb.String()
	}

	if len(result.Renamed) > 0 {
		sb.WriteString(fmt.Sprintf("Renamed (%d):\n", len(result.Renamed)))
		keys := make([]string, 0, len(result.Renamed))
		for k := range result.Renamed {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, old := range keys {
			sb.WriteString(fmt.Sprintf("  %s -> %s\n", old, result.Renamed[old]))
		}
	}

	if len(result.NotFound) > 0 {
		nf := make([]string, len(result.NotFound))
		copy(nf, result.NotFound)
		sort.Strings(nf)
		sb.WriteString(fmt.Sprintf("Not found (%d):\n", len(nf)))
		for _, k := range nf {
			sb.WriteString(fmt.Sprintf("  %s\n", k))
		}
	}

	return sb.String()
}

func renderRenameJSON(result diff.RenameResult) string {
	type payload struct {
		Renamed  map[string]string `json:"renamed"`
		NotFound []string          `json:"not_found"`
	}

	nf := result.NotFound
	if nf == nil {
		nf = []string{}
	}
	sort.Strings(nf)

	p := payload{
		Renamed:  result.Renamed,
		NotFound: nf,
	}
	if p.Renamed == nil {
		p.Renamed = map[string]string{}
	}

	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return `{"error": "failed to marshal rename result"}`
	}
	return string(b) + "\n"
}
