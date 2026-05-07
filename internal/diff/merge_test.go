package diff

import (
	"sort"
	"testing"
)

func TestMerge_NoOverrides(t *testing.T) {
	base := map[string]string{"HOST": "localhost", "PORT": "5432"}
	result := Merge(base)

	if len(result.Merged) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result.Merged))
	}
	if result.Merged["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %s", result.Merged["HOST"])
	}
	if len(result.Added) != 0 {
		t.Errorf("expected no added keys, got %v", result.Added)
	}
	if len(result.Overridden) != 0 {
		t.Errorf("expected no overridden keys, got %v", result.Overridden)
	}
}

func TestMerge_AddsNewKeys(t *testing.T) {
	base := map[string]string{"HOST": "localhost"}
	ov := map[string]string{"PORT": "3000"}
	result := Merge(base, ov)

	if result.Merged["PORT"] != "3000" {
		t.Errorf("expected PORT=3000, got %s", result.Merged["PORT"])
	}
	if len(result.Added) != 1 || result.Added[0] != "PORT" {
		t.Errorf("expected Added=[PORT], got %v", result.Added)
	}
}

func TestMerge_OverridesExistingKey(t *testing.T) {
	base := map[string]string{"HOST": "localhost", "PORT": "5432"}
	ov := map[string]string{"PORT": "9999"}
	result := Merge(base, ov)

	if result.Merged["PORT"] != "9999" {
		t.Errorf("expected PORT=9999, got %s", result.Merged["PORT"])
	}
	entry, ok := result.Overridden["PORT"]
	if !ok {
		t.Fatal("expected PORT in Overridden map")
	}
	if entry.Base != "5432" || entry.Override != "9999" {
		t.Errorf("unexpected override entry: %+v", entry)
	}
}

func TestMerge_MultipleOverrides(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	ov1 := map[string]string{"B": "20", "C": "30"}
	ov2 := map[string]string{"C": "300", "D": "400"}
	result := Merge(base, ov1, ov2)

	if result.Merged["B"] != "20" {
		t.Errorf("expected B=20, got %s", result.Merged["B"])
	}
	if result.Merged["C"] != "300" {
		t.Errorf("expected C=300, got %s", result.Merged["C"])
	}
	if result.Merged["D"] != "400" {
		t.Errorf("expected D=400, got %s", result.Merged["D"])
	}

	sort.Strings(result.Added)
	if len(result.Added) != 2 || result.Added[0] != "C" || result.Added[1] != "D" {
		t.Errorf("expected Added=[C D], got %v", result.Added)
	}
}

func TestMerge_BaseNotMutated(t *testing.T) {
	base := map[string]string{"KEY": "original"}
	ov := map[string]string{"KEY": "changed"}
	Merge(base, ov)

	if base["KEY"] != "original" {
		t.Errorf("base map was mutated, got %s", base["KEY"])
	}
}
