package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeScore(pct float64, total, matched, missing, extra, mismatched int) diff.Score {
	return diff.Score{
		Percent:    pct,
		Total:      total,
		Matched:    matched,
		Missing:    missing,
		Extra:      extra,
		Mismatched: mismatched,
	}
}

func TestRenderScore_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	s := makeScore(75.0, 4, 3, 1, 0, 0)
	if err := RenderScore(&buf, s, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "75.0%") {
		t.Errorf("expected percent in output, got: %s", out)
	}
	if !strings.Contains(out, "Compatibility Score") {
		t.Errorf("expected header in output, got: %s", out)
	}
}

func TestRenderScore_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	s := makeScore(50.0, 4, 2, 1, 1, 0)
	if err := RenderScore(&buf, s, "json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["percent"] != 50.0 {
		t.Errorf("expected percent 50.0, got %v", out["percent"])
	}
	if int(out["total_keys"].(float64)) != 4 {
		t.Errorf("expected total_keys 4, got %v", out["total_keys"])
	}
	if int(out["matched"].(float64)) != 2 {
		t.Errorf("expected matched 2, got %v", out["matched"])
	}
}

func TestRenderScore_DefaultsToText(t *testing.T) {
	var buf bytes.Buffer
	s := makeScore(100.0, 2, 2, 0, 0, 0)
	if err := RenderScore(&buf, s, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "100.0%") {
		t.Errorf("expected 100.0%% in default text output")
	}
}
