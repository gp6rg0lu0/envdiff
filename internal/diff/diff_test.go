package diff_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/diff"
)

func TestCompare_NoDifferences(t *testing.T) {
	base := map[string]string{"KEY": "val", "PORT": "8080"}
	target := map[string]string{"KEY": "val", "PORT": "8080"}

	result := diff.Compare(base, target)
	if result.HasDifferences() {
		t.Errorf("expected no differences, got %+v", result)
	}
}

func TestCompare_MissingKey(t *testing.T) {
	base := map[string]string{"KEY": "val", "SECRET": "s"}
	target := map[string]string{"KEY": "val"}

	result := diff.Compare(base, target)
	if len(result.Missing) != 1 || result.Missing[0] != "SECRET" {
		t.Errorf("expected SECRET in missing, got %v", result.Missing)
	}
}

func TestCompare_ExtraKey(t *testing.T) {
	base := map[string]string{"KEY": "val"}
	target := map[string]string{"KEY": "val", "EXTRA": "x"}

	result := diff.Compare(base, target)
	if len(result.Extra) != 1 || result.Extra[0] != "EXTRA" {
		t.Errorf("expected EXTRA in extra, got %v", result.Extra)
	}
}

func TestCompare_MismatchedValue(t *testing.T) {
	base := map[string]string{"PORT": "8080"}
	target := map[string]string{"PORT": "9090"}

	result := diff.Compare(base, target)
	pair, ok := result.Mismatched["PORT"]
	if !ok {
		t.Fatal("expected PORT in mismatched")
	}
	if pair[0] != "8080" || pair[1] != "9090" {
		t.Errorf("unexpected mismatch values: %v", pair)
	}
}

func TestCompare_AllDifferenceTypes(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2", "C": "3"}
	target := map[string]string{"A": "1", "B": "changed", "D": "4"}

	result := diff.Compare(base, target)
	if len(result.Missing) != 1 {
		t.Errorf("expected 1 missing, got %d", len(result.Missing))
	}
	if len(result.Extra) != 1 {
		t.Errorf("expected 1 extra, got %d", len(result.Extra))
	}
	if len(result.Mismatched) != 1 {
		t.Errorf("expected 1 mismatched, got %d", len(result.Mismatched))
	}
}
