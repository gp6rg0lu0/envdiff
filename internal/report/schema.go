package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// RenderSchema renders a SchemaResult in the requested format.
func RenderSchema(res diff.SchemaResult, format string) string {
	switch strings.ToLower(format) {
	case "json":
		return renderSchemaJSON(res)
	default:
		return renderSchemaText(res)
	}
}

func renderSchemaText(res diff.SchemaResult) string {
	var sb strings.Builder

	if res.Valid {
		sb.WriteString("schema: all rules passed\n")
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("schema: %d violation(s) found\n", len(res.Violations)))
	for _, v := range res.Violations {
		icon := "⚠"
		if v.Rule == "required" {
			icon = "✗"
		}
		sb.WriteString(fmt.Sprintf("  %s [%s] %s\n", icon, v.Key, v.Message))
	}
	return sb.String()
}

func renderSchemaJSON(res diff.SchemaResult) string {
	type jsonViolation struct {
		Key     string `json:"key"`
		Rule    string `json:"rule"`
		Message string `json:"message"`
	}
	type jsonOutput struct {
		Valid      bool            `json:"valid"`
		Violations []jsonViolation `json:"violations"`
	}

	out := jsonOutput{
		Valid:      res.Valid,
		Violations: make([]jsonViolation, 0, len(res.Violations)),
	}
	for _, v := range res.Violations {
		out.Violations = append(out.Violations, jsonViolation{
			Key:     v.Key,
			Rule:    v.Rule,
			Message: v.Message,
		})
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return `{"error":"failed to marshal schema result"}`
	}
	return string(b) + "\n"
}
