package diff

import (
	"strings"
	"testing"
)

func makePatchResult() Result {
	return Result{
		Missing: map[string]string{
			"DB_HOST": "localhost",
		},
		Extra: map[string]string{
			"OLD_KEY": "deprecated",
		},
		Mismatched: map[string][2]string{
			"APP_ENV": {"development", "production"},
		},
	}
}

func TestPatch_ReturnsAllOps(t *testing.T) {
	result := makePatchResult()
	entries := Patch(result)

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	ops := map[string]bool{}
	for _, e := range entries {
		ops[e.Op] = true
	}
	if !ops["add"] || !ops["remove"] || !ops["change"] {
		t.Errorf("expected add, remove, and change ops; got %v", ops)
	}
}

func TestPatch_AddEntry(t *testing.T) {
	result := Result{
		Missing:    map[string]string{"NEW_KEY": "value1"},
		Extra:      map[string]string{},
		Mismatched: map[string][2]string{},
	}
	entries := Patch(result)
	if len(entries) != 1 || entries[0].Op != "add" || entries[0].Key != "NEW_KEY" {
		t.Errorf("unexpected entry: %+v", entries)
	}
}

func TestPatch_SortedByKey(t *testing.T) {
	result := Result{
		Missing:    map[string]string{"Z_KEY": "z", "A_KEY": "a"},
		Extra:      map[string]string{},
		Mismatched: map[string][2]string{},
	}
	entries := Patch(result)
	if entries[0].Key != "A_KEY" || entries[1].Key != "Z_KEY" {
		t.Errorf("entries not sorted: %v, %v", entries[0].Key, entries[1].Key)
	}
}

func TestFormatPatch_AddLine(t *testing.T) {
	entries := []PatchEntry{{Key: "FOO", NewVal: "bar", Op: "add"}}
	out := FormatPatch(entries)
	if !strings.Contains(out, "+ FOO=bar") {
		t.Errorf("expected add line, got: %s", out)
	}
}

func TestFormatPatch_RemoveLine(t *testing.T) {
	entries := []PatchEntry{{Key: "OLD", OldVal: "gone", Op: "remove"}}
	out := FormatPatch(entries)
	if !strings.Contains(out, "- OLD=gone") {
		t.Errorf("expected remove line, got: %s", out)
	}
}

func TestFormatPatch_ChangeLine(t *testing.T) {
	entries := []PatchEntry{{Key: "ENV", OldVal: "dev", NewVal: "prod", Op: "change"}}
	out := FormatPatch(entries)
	if !strings.Contains(out, "~ ENV: dev -> prod") {
		t.Errorf("expected change line, got: %s", out)
	}
}

func TestFormatPatch_Empty(t *testing.T) {
	out := FormatPatch(nil)
	if out != "" {
		t.Errorf("expected empty string, got: %q", out)
	}
}
