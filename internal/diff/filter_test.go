package diff

import (
	"testing"
)

func baseResult() Result {
	return Result{
		Missing: []string{"DB_HOST"},
		Extra:   []string{"OLD_KEY"},
		Mismatched: []MismatchedEntry{
			{Key: "APP_ENV", BaseValue: "prod", OtherValue: "dev"},
		},
	}
}

func TestFilter_AllEnabled(t *testing.T) {
	r := Filter(baseResult(), DefaultFilterOptions())
	if len(r.Missing) != 1 || len(r.Extra) != 1 || len(r.Mismatched) != 1 {
		t.Errorf("expected all entries preserved, got missing=%d extra=%d mismatched=%d",
			len(r.Missing), len(r.Extra), len(r.Mismatched))
	}
}

func TestFilter_OnlyMissing(t *testing.T) {
	opts := FilterOptions{ShowMissing: true}
	r := Filter(baseResult(), opts)
	if len(r.Missing) != 1 {
		t.Errorf("expected 1 missing, got %d", len(r.Missing))
	}
	if len(r.Extra) != 0 || len(r.Mismatched) != 0 {
		t.Error("expected extra and mismatched to be empty")
	}
}

func TestFilter_OnlyExtra(t *testing.T) {
	opts := FilterOptions{ShowExtra: true}
	r := Filter(baseResult(), opts)
	if len(r.Extra) != 1 {
		t.Errorf("expected 1 extra, got %d", len(r.Extra))
	}
	if len(r.Missing) != 0 || len(r.Mismatched) != 0 {
		t.Error("expected missing and mismatched to be empty")
	}
}

func TestFilter_OnlyMismatched(t *testing.T) {
	opts := FilterOptions{ShowMismatched: true}
	r := Filter(baseResult(), opts)
	if len(r.Mismatched) != 1 {
		t.Errorf("expected 1 mismatched, got %d", len(r.Mismatched))
	}
	if len(r.Missing) != 0 || len(r.Extra) != 0 {
		t.Error("expected missing and extra to be empty")
	}
}

func TestFilter_NoneEnabled(t *testing.T) {
	r := Filter(baseResult(), FilterOptions{})
	if !IsEmpty(r) {
		t.Error("expected empty result when no options enabled")
	}
}

func TestIsEmpty_True(t *testing.T) {
	r := Result{
		Missing:    []string{},
		Extra:      []string{},
		Mismatched: []MismatchedEntry{},
	}
	if !IsEmpty(r) {
		t.Error("expected IsEmpty to return true")
	}
}

func TestIsEmpty_False(t *testing.T) {
	if IsEmpty(baseResult()) {
		t.Error("expected IsEmpty to return false")
	}
}
