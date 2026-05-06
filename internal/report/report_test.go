package report_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/report"
)

func makeResult() diff.Result {
	return diff.Result{
		Missing: []string{"DB_HOST"},
		Extra:   []string{"NEW_KEY"},
		Mismatched: []diff.Mismatch{
			{Key: "APP_ENV", BaseValue: "production", CompareValue: "staging"},
		},
	}
}

func TestRenderText_NoDifferences(t *testing.T) {
	var buf bytes.Buffer
	err := report.Render(&buf, diff.Result{}, "a.env", "b.env", report.Options{Format: report.FormatText, NoColor: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected clean message, got: %s", buf.String())
	}
}

func TestRenderText_ShowsDifferences(t *testing.T) {
	var buf bytes.Buffer
	err := report.Render(&buf, makeResult(), "base.env", "compare.env", report.Options{Format: report.FormatText, NoColor: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "MISSING") {
		t.Errorf("expected MISSING in output, got: %s", out)
	}
	if !strings.Contains(out, "EXTRA") {
		t.Errorf("expected EXTRA in output, got: %s", out)
	}
	if !strings.Contains(out, "MISMATCH") {
		t.Errorf("expected MISMATCH in output, got: %s", out)
	}
}

func TestRenderJSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	err := report.Render(&buf, makeResult(), "base.env", "compare.env", report.Options{Format: report.FormatJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if out["clean"] != false {
		t.Errorf("expected clean=false")
	}
	if out["base_file"] != "base.env" {
		t.Errorf("expected base_file=base.env")
	}
}

func TestRenderJSON_CleanResult(t *testing.T) {
	var buf bytes.Buffer
	err := report.Render(&buf, diff.Result{}, "a.env", "b.env", report.Options{Format: report.FormatJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["clean"] != true {
		t.Errorf("expected clean=true for empty result")
	}
}
