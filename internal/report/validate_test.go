package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeValidationResult(issues ...diff.ValidationIssue) diff.ValidationResult {
	return diff.ValidationResult{Issues: issues}
}

func TestRenderValidationText_Valid(t *testing.T) {
	var buf bytes.Buffer
	err := RenderValidation(&buf, makeValidationResult(), "text")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "passed") {
		t.Errorf("expected 'passed' in output, got: %s", buf.String())
	}
}

func TestRenderValidationText_WithIssues(t *testing.T) {
	var buf bytes.Buffer
	issues := []diff.ValidationIssue{
		{Key: "HOST", Message: "value is empty"},
		{Key: "PORT", Message: "required key is missing"},
	}
	err := RenderValidation(&buf, makeValidationResult(issues...), "text")
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "failed") {
		t.Errorf("expected 'failed' in output")
	}
	if !strings.Contains(out, "HOST") || !strings.Contains(out, "PORT") {
		t.Errorf("expected issue keys in output, got: %s", out)
	}
}

func TestRenderValidationJSON_Valid(t *testing.T) {
	var buf bytes.Buffer
	err := RenderValidation(&buf, makeValidationResult(), "json")
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["valid"] != true {
		t.Errorf("expected valid=true, got %v", out["valid"])
	}
	if out["issue_count"].(float64) != 0 {
		t.Errorf("expected issue_count=0")
	}
}

func TestRenderValidationJSON_WithIssues(t *testing.T) {
	var buf bytes.Buffer
	issues := []diff.ValidationIssue{{Key: "SECRET", Message: "forbidden key is present"}}
	err := RenderValidation(&buf, makeValidationResult(issues...), "json")
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["valid"] != false {
		t.Errorf("expected valid=false")
	}
	if out["issue_count"].(float64) != 1 {
		t.Errorf("expected issue_count=1")
	}
}

func TestRenderValidation_DefaultsToText(t *testing.T) {
	var buf bytes.Buffer
	err := RenderValidation(&buf, makeValidationResult(), "")
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(buf.String()) == "" {
		t.Error("expected non-empty output for default format")
	}
}
