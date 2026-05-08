package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// RenderPatch renders patch entries in the specified format.
func RenderPatch(entries []diff.PatchEntry, format string) string {
	switch strings.ToLower(format) {
	case "json":
		return renderPatchJSON(entries)
	default:
		return renderPatchText(entries)
	}
}

func renderPatchText(entries []diff.PatchEntry) string {
	if len(entries) == 0 {
		return "No changes.\n"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Patch (%d change(s)):\n", len(entries)))

	for _, e := range entries {
		switch e.Op {
		case "add":
			sb.WriteString(fmt.Sprintf("  + %s=%s\n", e.Key, e.NewVal))
		case "remove":
			sb.WriteString(fmt.Sprintf("  - %s=%s\n", e.Key, e.OldVal))
		case "change":
			sb.WriteString(fmt.Sprintf("  ~ %s: %q -> %q\n", e.Key, e.OldVal, e.NewVal))
		}
	}

	return sb.String()
}

type patchEntryJSON struct {
	Key    string `json:"key"`
	Op     string `json:"op"`
	OldVal string `json:"old_value,omitempty"`
	NewVal string `json:"new_value,omitempty"`
}

func renderPatchJSON(entries []diff.PatchEntry) string {
	records := make([]patchEntryJSON, 0, len(entries))
	for _, e := range entries {
		records = append(records, patchEntryJSON{
			Key:    e.Key,
			Op:     e.Op,
			OldVal: e.OldVal,
			NewVal: e.NewVal,
		})
	}

	payload := map[string]interface{}{
		"changes": records,
		"count":   len(records),
	}

	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return `{"error": "failed to marshal patch"}`
	}
	return string(b) + "\n"
}
