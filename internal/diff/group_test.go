package diff

import (
	"testing"
)

func makeGroupResult() Result {
	return Result{
		Entries: []Entry{
			{Key: "ZEBRA", Status: StatusMissing},
			{Key: "APPLE", Status: StatusExtra},
			{Key: "MANGO", Status: StatusMismatch, BaseValue: "a", OtherValue: "b"},
			{Key: "BANANA", Status: StatusMissing},
			{Key: "CHERRY", Status: StatusMismatch, BaseValue: "x", OtherValue: "y"},
		},
	}
}

func TestGroup_SplitsByStatus(t *testing.T) {
	r := makeGroupResult()
	g := Group(r)

	if len(g.Missing) != 2 {
		t.Errorf("expected 2 missing, got %d", len(g.Missing))
	}
	if len(g.Extra) != 1 {
		t.Errorf("expected 1 extra, got %d", len(g.Extra))
	}
	if len(g.Mismatched) != 2 {
		t.Errorf("expected 2 mismatched, got %d", len(g.Mismatched))
	}
}

func TestGroup_MissingSortedByKey(t *testing.T) {
	r := makeGroupResult()
	g := Group(r)

	if g.Missing[0].Key != "BANANA" || g.Missing[1].Key != "ZEBRA" {
		t.Errorf("missing entries not sorted: %v", g.Missing)
	}
}

func TestGroup_MismatchedSortedByKey(t *testing.T) {
	r := makeGroupResult()
	g := Group(r)

	if g.Mismatched[0].Key != "CHERRY" || g.Mismatched[1].Key != "MANGO" {
		t.Errorf("mismatched entries not sorted: %v", g.Mismatched)
	}
}

func TestGroup_EmptyResult(t *testing.T) {
	g := Group(Result{})
	if g.Total() != 0 {
		t.Errorf("expected total 0, got %d", g.Total())
	}
}

func TestGroup_Total(t *testing.T) {
	r := makeGroupResult()
	g := Group(r)
	if g.Total() != 5 {
		t.Errorf("expected total 5, got %d", g.Total())
	}
}

func TestGroup_Keys(t *testing.T) {
	r := makeGroupResult()
	g := Group(r)
	keys := g.Keys()

	expected := []string{"APPLE", "BANANA", "CHERRY", "MANGO", "ZEBRA"}
	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(keys))
	}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("keys[%d]: expected %q, got %q", i, expected[i], k)
		}
	}
}

func TestGroup_NoDuplicateKeys(t *testing.T) {
	r := Result{
		Entries: []Entry{
			{Key: "FOO", Status: StatusMissing},
			{Key: "FOO", Status: StatusExtra},
		},
	}
	g := Group(r)
	keys := g.Keys()
	if len(keys) != 1 {
		t.Errorf("expected 1 unique key, got %d: %v", len(keys), keys)
	}
}
