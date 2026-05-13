package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// RenderLint formats a LintResult as either text or JSON.
func RenderLint(result diff.LintResult, format string) string {
	switch strings.ToLower(format) {
	case "json":
		return renderLintJSON(result)
	default:
		return renderLintText(result)
	}
}

func renderLintText(result diff.LintResult) string {
	if result.Valid {
		return "lint: no issues found\n"
	}

	// Sort issues by key then severity for deterministic output.
	issues := make([]diff.LintIssue, len(result.Issues))
	copy(issues, result.Issues)
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Key != issues[j].Key {
			return issues[i].Key < issues[j].Key
		}
		return issues[i].Severity < issues[j].Severity
	})

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("lint: %d issue(s) found\n", len(issues)))
	for _, issue := range issues {
		icon := lintIcon(issue.Severity)
		sb.WriteString(fmt.Sprintf("  %s [%s] %s: %s\n", icon, issue.Severity, issue.Key, issue.Message))
	}
	return sb.String()
}

func lintIcon(s diff.LintSeverity) string {
	switch s {
	case diff.LintError:
		return "✗"
	case diff.LintWarning:
		return "⚠"
	default:
		return "ℹ"
	}
}

func renderLintJSON(result diff.LintResult) string {
	type jsonIssue struct {
		Key      string `json:"key"`
		Message  string `json:"message"`
		Severity string `json:"severity"`
	}
	type jsonOutput struct {
		Valid  bool        `json:"valid"`
		Count  int         `json:"issue_count"`
		Issues []jsonIssue `json:"issues"`
	}

	out := jsonOutput{
		Valid:  result.Valid,
		Count:  len(result.Issues),
		Issues: make([]jsonIssue, 0, len(result.Issues)),
	}
	for _, iss := range result.Issues {
		out.Issues = append(out.Issues, jsonIssue{
			Key:      iss.Key,
			Message:  iss.Message,
			Severity: string(iss.Severity),
		})
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	return string(b) + "\n"
}
