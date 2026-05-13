package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeSchemaResult(valid bool, violations ...diff.SchemaViolation) diff.SchemaResult {
	return diff.SchemaResult{Valid: valid, Violations: violations}
}

func TestRenderSchemaText_Valid(t *testing.T) {
	res := makeSchemaResult(true)
	out := RenderSchema(res, "text")
	if !strings.Contains(out, "all rules passed") {
		t.Errorf("expected pass message, got: %q", out)
	}
}

func TestRenderSchemaText_WithViolations(t *testing.T) {
	res := makeSchemaResult(false,
		diff.SchemaViolation{Key: "PORT", Rule: "type", Message: `key "PORT": expected integer, got "abc"`},
		diff.SchemaViolation{Key: "SECRET", Rule: "required", Message: `key "SECRET" is required but missing`},
	)
	out := RenderSchema(res, "text")
	if !strings.Contains(out, "2 violation(s)") {
		t.Errorf("expected violation count, got: %q", out)
	}
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output, got: %q", out)
	}
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected SECRET in output, got: %q", out)
	}
}

func TestRenderSchemaText_RequiredIconDiffers(t *testing.T) {
	res := makeSchemaResult(false,
		diff.SchemaViolation{Key: "MISSING", Rule: "required", Message: "missing"},
	)
	out := RenderSchema(res, "text")
	if !strings.Contains(out, "✗") {
		t.Errorf("expected ✗ icon for required violation, got: %q", out)
	}
}

func TestRenderSchemaJSON_Valid(t *testing.T) {
	res := makeSchemaResult(true)
	out := RenderSchema(res, "json")
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(out), &data); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if data["valid"] != true {
		t.Errorf("expected valid=true")
	}
	violations, ok := data["violations"].([]interface{})
	if !ok || len(violations) != 0 {
		t.Errorf("expected empty violations array")
	}
}

func TestRenderSchemaJSON_WithViolations(t *testing.T) {
	res := makeSchemaResult(false,
		diff.SchemaViolation{Key: "PORT", Rule: "type", Message: "bad type"},
	)
	out := RenderSchema(res, "json")
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(out), &data); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if data["valid"] != false {
		t.Errorf("expected valid=false")
	}
	violations := data["violations"].([]interface{})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	v := violations[0].(map[string]interface{})
	if v["key"] != "PORT" {
		t.Errorf("expected key PORT, got %v", v["key"])
	}
}

func TestRenderSchema_DefaultsToText(t *testing.T) {
	res := makeSchemaResult(true)
	out := RenderSchema(res, "")
	if !strings.Contains(out, "all rules passed") {
		t.Errorf("expected text fallback, got: %q", out)
	}
}
