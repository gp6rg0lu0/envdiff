package diff

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
)

func makeSortResult() parser.Result {
	return parser.Result{
		BaseFile:    "a.env",
		CompareFile: "b.env",
		Differences: []parser.DiffEntry{
			{Key: "ZEBRA", Extra: true},
			{Key: "APPLE", Missing: true},
			{Key: "MANGO", Mismatched: true, BaseValue: "x", CompareValue: "y"},
			{Key: "BANANA", Missing: true},
			{Key: "CHERRY", Extra: true},
		},
	}
}

func TestSortResult_ByKey(t *testing.T) {
	r := SortResult(makeSortResult(), SortByKey)
	keys := make([]string, len(r.Differences))
	for i, d := range r.Differences {
		keys[i] = d.Key
	}
	expected := []string{"APPLE", "BANANA", "CHERRY", "MANGO", "ZEBRA"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("position %d: got %q, want %q", i, keys[i], k)
		}
	}
}

func TestSortResult_ByStatus(t *testing.T) {
	r := SortResult(makeSortResult(), SortByStatus)
	// missing (rank 0) should come before extra (rank 1) before mismatched (rank 2)
	if !r.Differences[0].Missing || !r.Differences[1].Missing {
		t.Error("expected first two entries to be Missing")
	}
	if !r.Differences[2].Extra || !r.Differences[3].Extra {
		t.Error("expected next two entries to be Extra")
	}
	if !r.Differences[4].Mismatched {
		t.Error("expected last entry to be Mismatched")
	}
}

func TestSortResult_ByStatusKey(t *testing.T) {
	r := SortResult(makeSortResult(), SortByStatusKey)
	// missing keys sorted: APPLE, BANANA
	if r.Differences[0].Key != "APPLE" {
		t.Errorf("expected APPLE first, got %q", r.Differences[0].Key)
	}
	if r.Differences[1].Key != "BANANA" {
		t.Errorf("expected BANANA second, got %q", r.Differences[1].Key)
	}
	// extra keys sorted: CHERRY, ZEBRA
	if r.Differences[2].Key != "CHERRY" {
		t.Errorf("expected CHERRY third, got %q", r.Differences[2].Key)
	}
	if r.Differences[3].Key != "ZEBRA" {
		t.Errorf("expected ZEBRA fourth, got %q", r.Differences[3].Key)
	}
}

func TestSortResult_DoesNotMutateOriginal(t *testing.T) {
	orig := makeSortResult()
	firstKey := orig.Differences[0].Key
	SortResult(orig, SortByKey)
	if orig.Differences[0].Key != firstKey {
		t.Error("SortResult mutated the original result")
	}
}

func TestSortResult_UnknownOrder(t *testing.T) {
	orig := makeSortResult()
	r := SortResult(orig, SortOrder("unknown"))
	for i, d := range r.Differences {
		if d.Key != orig.Differences[i].Key {
			t.Errorf("expected unchanged order at %d", i)
		}
	}
}
