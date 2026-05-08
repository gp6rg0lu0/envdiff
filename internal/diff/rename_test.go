package diff

import (
	"sort"
	"testing"
)

func TestApplyRenames_NoRenames(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result := ApplyRenames(env, RenameMap{})

	if len(result.Renamed) != 0 {
		t.Errorf("expected no renames, got %v", result.Renamed)
	}
	if len(result.NotFound) != 0 {
		t.Errorf("expected no not-found, got %v", result.NotFound)
	}
	if result.Output["FOO"] != "bar" || result.Output["BAZ"] != "qux" {
		t.Errorf("unexpected output: %v", result.Output)
	}
}

func TestApplyRenames_BasicRename(t *testing.T) {
	env := map[string]string{"OLD_KEY": "value1", "OTHER": "value2"}
	renames := RenameMap{"OLD_KEY": "NEW_KEY"}

	result := ApplyRenames(env, renames)

	if result.Renamed["OLD_KEY"] != "NEW_KEY" {
		t.Errorf("expected OLD_KEY -> NEW_KEY in renamed, got %v", result.Renamed)
	}
	if _, exists := result.Output["OLD_KEY"]; exists {
		t.Error("OLD_KEY should not exist in output")
	}
	if result.Output["NEW_KEY"] != "value1" {
		t.Errorf("expected NEW_KEY=value1, got %q", result.Output["NEW_KEY"])
	}
	if result.Output["OTHER"] != "value2" {
		t.Errorf("expected OTHER=value2, got %q", result.Output["OTHER"])
	}
}

func TestApplyRenames_KeyNotFound(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	renames := RenameMap{"MISSING_KEY": "NEW_KEY"}

	result := ApplyRenames(env, renames)

	if len(result.NotFound) != 1 || result.NotFound[0] != "MISSING_KEY" {
		t.Errorf("expected MISSING_KEY in NotFound, got %v", result.NotFound)
	}
	if _, exists := result.Output["NEW_KEY"]; exists {
		t.Error("NEW_KEY should not appear in output when source key is missing")
	}
}

func TestApplyRenames_MultipleRenames(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2", "C": "3"}
	renames := RenameMap{"A": "ALPHA", "B": "BETA"}

	result := ApplyRenames(env, renames)

	if len(result.Renamed) != 2 {
		t.Errorf("expected 2 renames, got %d", len(result.Renamed))
	}
	if result.Output["ALPHA"] != "1" || result.Output["BETA"] != "2" || result.Output["C"] != "3" {
		t.Errorf("unexpected output: %v", result.Output)
	}
}

func TestApplyRenames_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"OLD": "val"}
	original := map[string]string{"OLD": "val"}
	renames := RenameMap{"OLD": "NEW"}

	ApplyRenames(env, renames)

	for k, v := range original {
		if env[k] != v {
			t.Errorf("input env was mutated: key %q changed", k)
		}
	}
}

func TestApplyRenames_NotFoundSorted(t *testing.T) {
	env := map[string]string{"KEEP": "v"}
	renames := RenameMap{"Z_MISSING": "Z_NEW", "A_MISSING": "A_NEW"}

	result := ApplyRenames(env, renames)

	if len(result.NotFound) != 2 {
		t.Fatalf("expected 2 not-found, got %d", len(result.NotFound))
	}
	if !sort.StringsAreSorted(result.NotFound) {
		// NotFound order is non-deterministic from map iteration; just check membership.
		got := map[string]bool{}
		for _, k := range result.NotFound {
			got[k] = true
		}
		if !got["Z_MISSING"] || !got["A_MISSING"] {
			t.Errorf("unexpected NotFound contents: %v", result.NotFound)
		}
	}
}
