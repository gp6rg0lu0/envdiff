package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// RenderOverlap renders an OverlapResult in the given format.
func RenderOverlap(r diff.OverlapResult, format string) string {
	switch format {
	case "json":
		return renderOverlapJSON(r)
	default:
		return renderOverlapText(r)
	}
}

func renderOverlapText(r diff.OverlapResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Overlap: %.1f%%\n", r.OverlapPercent))
	sb.WriteString(fmt.Sprintf("  Shared keys    : %d\n", len(r.SharedKeys)))
	sb.WriteString(fmt.Sprintf("  Unique to base : %d\n", len(r.UniqueToBase)))
	sb.WriteString(fmt.Sprintf("  Unique to cmp  : %d\n", len(r.UniqueToCompare)))

	if len(r.SharedKeys) > 0 {
		sb.WriteString("  [shared]\n")
		for _, k := range r.SharedKeys {
			sb.WriteString(fmt.Sprintf("    = %s\n", k))
		}
	}
	if len(r.UniqueToBase) > 0 {
		sb.WriteString("  [only in base]\n")
		for _, k := range r.UniqueToBase {
			sb.WriteString(fmt.Sprintf("    - %s\n", k))
		}
	}
	if len(r.UniqueToCompare) > 0 {
		sb.WriteString("  [only in compare]\n")
		for _, k := range r.UniqueToCompare {
			sb.WriteString(fmt.Sprintf("    + %s\n", k))
		}
	}
	return sb.String()
}

func renderOverlapJSON(r diff.OverlapResult) string {
	payload := map[string]interface{}{
		"overlap_percent":   r.OverlapPercent,
		"shared_keys":       r.SharedKeys,
		"unique_to_base":    r.UniqueToBase,
		"unique_to_compare": r.UniqueToCompare,
	}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return `{"error":"failed to marshal overlap result"}`
	}
	return string(b) + "\n"
}
