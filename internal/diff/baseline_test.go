package diff

import (
	"testing"
)

func TestNewBaseline_SortedKeys(t *testing.T) {
	env := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	b := NewBaseline("test", env)

	if b.Name != "test" {
		t.Errorf("expected name 'test', got %q", b.Name)
	}
	if len(b.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(b.Entries))
	}
	if b.Entries[0].Key != "A_KEY" || b.Entries[1].Key != "M_KEY" || b.Entries[2].Key != "Z_KEY" {
		t.Errorf("entries not sorted: %v", b.Entries)
	}
}

func TestBaseline_ToMap_RoundTrip(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := NewBaseline("snap", env)
	result := b.ToMap()

	for k, v := range env {
		if result[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, result[k])
		}
	}
}

func TestDiffBaseline_NoChanges(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "5432"}
	b := NewBaseline("prod", env)
	r := DiffBaseline(b, env)

	if len(r.Missing) != 0 || len(r.Extra) != 0 || len(r.Mismatched) != 0 {
		t.Errorf("expected no differences, got %+v", r)
	}
}

func TestDiffBaseline_DetectsNewKey(t *testing.T) {
	base := map[string]string{"HOST": "localhost"}
	current := map[string]string{"HOST": "localhost", "NEW_KEY": "value"}
	b := NewBaseline("snap", base)
	r := DiffBaseline(b, current)

	if len(r.Extra) != 1 || r.Extra[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY as extra, got %v", r.Extra)
	}
}

func TestDiffBaseline_DetectsRemovedKey(t *testing.T) {
	base := map[string]string{"HOST": "localhost", "OLD_KEY": "legacy"}
	current := map[string]string{"HOST": "localhost"}
	b := NewBaseline("snap", base)
	r := DiffBaseline(b, current)

	if len(r.Missing) != 1 || r.Missing[0] != "OLD_KEY" {
		t.Errorf("expected OLD_KEY as missing, got %v", r.Missing)
	}
}

func TestDiffBaseline_DetectsValueChange(t *testing.T) {
	base := map[string]string{"DB_URL": "postgres://old"}
	current := map[string]string{"DB_URL": "postgres://new"}
	b := NewBaseline("snap", base)
	r := DiffBaseline(b, current)

	if len(r.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(r.Mismatched))
	}
	if r.Mismatched["DB_URL"].Base != "postgres://old" || r.Mismatched["DB_URL"].Compare != "postgres://new" {
		t.Errorf("unexpected mismatch values: %+v", r.Mismatched["DB_URL"])
	}
}
