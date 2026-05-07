package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// RenderValidation writes a validation result to w in the requested format.
func RenderValidation(w io.Writer, result diff.ValidationResult, format string) error {
	switch strings.ToLower(format) {
	case "json":
		return renderValidationJSON(w, result)
	default:
		return renderValidationText(w, result)
	}
}

func renderValidationText(w io.Writer, result diff.ValidationResult) error {
	if result.Valid() {
		_, err := fmt.Fprintln(w, "validation passed: no issues found")
		return err
	}
	_, err := fmt.Fprintf(w, "validation failed: %d issue(s) found\n", len(result.Issues))
	if err != nil {
		return err
	}
	for _, issue := range result.Issues {
		_, err = fmt.Fprintf(w, "  [INVALID] %s\n", issue)
		if err != nil {
			return err
		}
	}
	return nil
}

func renderValidationJSON(w io.Writer, result diff.ValidationResult) error {
	type jsonIssue struct {
		Key     string `json:"key"`
		Message string `json:"message"`
	}
	type jsonOutput struct {
		Valid  bool        `json:"valid"`
		Count  int         `json:"issue_count"`
		Issues []jsonIssue `json:"issues"`
	}

	issues := make([]jsonIssue, 0, len(result.Issues))
	for _, iss := range result.Issues {
		issues = append(issues, jsonIssue{Key: iss.Key, Message: iss.Message})
	}
	out := jsonOutput{
		Valid:  result.Valid(),
		Count:  len(result.Issues),
		Issues: issues,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
