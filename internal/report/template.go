package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// RenderTemplate renders a TemplateResult in the given format.
func RenderTemplate(res diff.TemplateResult, format string) string {
	switch strings.ToLower(format) {
	case "json":
		return renderTemplateJSON(res)
	default:
		return renderTemplateText(res)
	}
}

func renderTemplateText(res diff.TemplateResult) string {
	if len(res.Vars) == 0 {
		return "No template variable references found.\n"
	}

	var sb strings.Builder
	sb.WriteString("Template Variables:\n")

	for _, v := range res.Vars {
		if v.Resolved {
			sb.WriteString(fmt.Sprintf("  ✔ %s => %q\n", v.Key, v.Value))
		} else {
			sb.WriteString(fmt.Sprintf("  ✘ %s (unresolved)\n", v.Key))
		}
	}

	if len(res.Unresolved) > 0 {
		sb.WriteString(fmt.Sprintf("\nUnresolved references (%d): %s\n",
			len(res.Unresolved),
			strings.Join(res.Unresolved, ", "),
		))
	}
	return sb.String()
}

func renderTemplateJSON(res diff.TemplateResult) string {
	type varEntry struct {
		Key      string `json:"key"`
		Resolved bool   `json:"resolved"`
		Value    string `json:"value,omitempty"`
	}
	type output struct {
		Vars       []varEntry `json:"vars"`
		Unresolved []string   `json:"unresolved"`
	}

	out := output{
		Vars:       make([]varEntry, 0, len(res.Vars)),
		Unresolved: res.Unresolved,
	}
	if out.Unresolved == nil {
		out.Unresolved = []string{}
	}
	for _, v := range res.Vars {
		out.Vars = append(out.Vars, varEntry{
			Key:      v.Key,
			Resolved: v.Resolved,
			Value:    v.Value,
		})
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	return string(b) + "\n"
}
