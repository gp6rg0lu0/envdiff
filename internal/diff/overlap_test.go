package diff

import (
	"testing"
)

func TestOverlap_IdenticalMaps(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	cmp := map[string]string{"A": "1", "B": "2"}
	r := Overlap(base, cmp)
	if len(r.SharedKeys) != 2 {
		t.Errorf("expected 2 shared keys, got %d", len(r.SharedKeys))
	}
	if len(r.UniqueToBase) != 0 || len(r.UniqueToCompare) != 0 {
		t.Error("expected no unique keys")
	}
	if r.OverlapPercent != 100.0 {
		t.Errorf("expected 100%% overlap, got %.2f", r.OverlapPercent)
	}
}

func TestOverlap_NoCommonKeys(t *testing.T) {
	base := map[string]string{"A": "1"}
	cmp := map[string]string{"B": "2"}
	r := Overlap(base, cmp)
	if len(r.SharedKeys) != 0 {
		t.Errorf("expected 0 shared keys, got %d", len(r.SharedKeys))
	}
	if len(r.UniqueToBase) != 1 || r.UniqueToBase[0] != "A" {
		t.Errorf("expected UniqueToBase=[A], got %v", r.UniqueToBase)
	}
	if len(r.UniqueToCompare) != 1 || r.UniqueToCompare[0] != "B" {
		t.Errorf("expected UniqueToCompare=[B], got %v", r.UniqueToCompare)
	}
	if r.OverlapPercent != 0.0 {
		t.Errorf("expected 0%% overlap, got %.2f", r.OverlapPercent)
	}
}

func TestOverlap_PartialOverlap(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2", "C": "3"}
	cmp := map[string]string{"B": "2", "C": "3", "D": "4"}
	r := Overlap(base, cmp)
	if len(r.SharedKeys) != 2 {
		t.Errorf("expected 2 shared keys, got %d", len(r.SharedKeys))
	}
	if len(r.UniqueToBase) != 1 || r.UniqueToBase[0] != "A" {
		t.Errorf("expected UniqueToBase=[A], got %v", r.UniqueToBase)
	}
	if len(r.UniqueToCompare) != 1 || r.UniqueToCompare[0] != "D" {
		t.Errorf("expected UniqueToCompare=[D], got %v", r.UniqueToCompare)
	}
	// union=4, shared=2 => 50%
	if r.OverlapPercent != 50.0 {
		t.Errorf("expected 50%% overlap, got %.2f", r.OverlapPercent)
	}
}

func TestOverlap_EmptyMaps(t *testing.T) {
	r := Overlap(map[string]string{}, map[string]string{})
	if r.OverlapPercent != 0.0 {
		t.Errorf("expected 0%% for empty maps, got %.2f", r.OverlapPercent)
	}
	if len(r.SharedKeys) != 0 {
		t.Error("expected no shared keys for empty maps")
	}
}

func TestOverlap_SortedOutput(t *testing.T) {
	base := map[string]string{"Z": "1", "A": "2", "M": "3"}
	cmp := map[string]string{"Z": "1", "A": "2", "M": "3"}
	r := Overlap(base, cmp)
	expected := []string{"A", "M", "Z"}
	for i, k := range r.SharedKeys {
		if k != expected[i] {
			t.Errorf("expected SharedKeys[%d]=%s, got %s", i, expected[i], k)
		}
	}
}
