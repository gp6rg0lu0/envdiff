package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeBaseline(name string, env map[string]string) diff.Baseline {
	return diff.NewBaseline(name, env)
}

func TestRenderBaselineText_NoChanges(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "5432"}
	b := makeBaseline("production", env)
	r := diff.DiffBaseline(b, env)

	var buf bytes.Buffer
	if err := RenderBaseline(&buf, b, r, "text", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "production") {
		t.Errorf("expected baseline name in output, got: %s", out)
	}
	if !strings.Contains(out, "No changes") {
		t.Errorf("expected 'No changes' message, got: %s", out)
	}
}

func TestRenderBaselineText_WithChanges(t *testing.T) {
	base := map[string]string{"HOST": "old", "REMOVED": "gone"}
	current := map[string]string{"HOST": "new", "ADDED": "here"}
	b := makeBaseline("staging", base)
	r := diff.DiffBaseline(b, current)

	var buf bytes.Buffer
	if err := RenderBaseline(&buf, b, r, "text", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "REMOVED") {
		t.Errorf("expected REMOVED key in output: %s", out)
	}
	if !strings.Contains(out, "ADDED") {
		t.Errorf("expected ADDED key in output: %s", out)
	}
	if !strings.Contains(out, "HOST") {
		t.Errorf("expected HOST change in output: %s", out)
	}
}

func TestRenderBaselineJSON_Structure(t *testing.T) {
	base := map[string]string{"KEY": "val"}
	current := map[string]string{"KEY": "changed"}
	b := makeBaseline("dev", base)
	r := diff.DiffBaseline(b, current)

	var buf bytes.Buffer
	if err := RenderBaseline(&buf, b, r, "json", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["baseline"] != "dev" {
		t.Errorf("expected baseline name 'dev', got %v", out["baseline"])
	}
	if _, ok := out["changes"]; !ok {
		t.Error("expected 'changes' key in JSON output")
	}
}

func TestRenderBaselineText_Redact(t *testing.T) {
	base := map[string]string{"SECRET_KEY": "topsecret"}
	current := map[string]string{"SECRET_KEY": "newsecret"}
	b := makeBaseline("prod", base)
	r := diff.DiffBaseline(b, current)

	var buf bytes.Buffer
	if err := RenderBaseline(&buf, b, r, "text", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "topsecret") || strings.Contains(out, "newsecret") {
		t.Errorf("sensitive values should be redacted, got: %s", out)
	}
}
