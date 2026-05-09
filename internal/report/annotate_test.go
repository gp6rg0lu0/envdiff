package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeAnnotations() []diff.Annotation {
	return []diff.Annotation{
		{Key: "DB_HOST", Level: diff.AnnotationError, Message: "key is missing from target environment"},
		{Key: "APP_NAME", Level: diff.AnnotationWarning, Message: "key exists only in target environment"},
		{Key: "LOG_LEVEL", Level: diff.AnnotationInfo, Message: "value differs between environments"},
	}
}

func TestRenderAnnotationsText_Empty(t *testing.T) {
	out := RenderAnnotations(nil, "text")
	if !strings.Contains(out, "No annotations") {
		t.Errorf("expected 'No annotations' message, got: %s", out)
	}
}

func TestRenderAnnotationsText_WithAnnotations(t *testing.T) {
	anns := makeAnnotations()
	out := RenderAnnotations(anns, "text")

	if !strings.Contains(out, "Annotations:") {
		t.Errorf("expected header 'Annotations:', got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected key DB_HOST in output")
	}
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected ERROR level in output")
	}
	if !strings.Contains(out, "WARNING") {
		t.Errorf("expected WARNING level in output")
	}
}

func TestRenderAnnotationsText_Icons(t *testing.T) {
	anns := makeAnnotations()
	out := RenderAnnotations(anns, "text")
	if !strings.Contains(out, "✖") {
		t.Errorf("expected error icon ✖")
	}
	if !strings.Contains(out, "⚠") {
		t.Errorf("expected warning icon ⚠")
	}
}

func TestRenderAnnotationsJSON_Structure(t *testing.T) {
	anns := makeAnnotations()
	out := RenderAnnotations(anns, "json")

	var result struct {
		Annotations []struct {
			Key     string `json:"key"`
			Level   string `json:"level"`
			Message string `json:"message"`
		} `json:"annotations"`
		Count int `json:"count"`
	}

	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Count != 3 {
		t.Errorf("expected count 3, got %d", result.Count)
	}
	if result.Annotations[0].Key != "DB_HOST" {
		t.Errorf("expected first key DB_HOST, got %s", result.Annotations[0].Key)
	}
}

func TestRenderAnnotationsJSON_Empty(t *testing.T) {
	out := RenderAnnotations(nil, "json")
	var result struct {
		Count int `json:"count"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Count != 0 {
		t.Errorf("expected count 0, got %d", result.Count)
	}
}

func TestRenderAnnotations_DefaultsToText(t *testing.T) {
	anns := makeAnnotations()
	out := RenderAnnotations(anns, "")
	if !strings.Contains(out, "Annotations:") {
		t.Errorf("expected text format as default")
	}
}
