package diff

import (
	"strings"
	"testing"
)

func TestSummarize_NoDifferences(t *testing.T) {
	r := Result{
		Missing:    []string{},
		Extra:      []string{},
		Mismatched: []MismatchedEntry{},
	}
	s := Summarize(r)
	if s.TotalCount != 0 {
		t.Errorf("expected TotalCount 0, got %d", s.TotalCount)
	}
	if s.HasDifferences() {
		t.Error("expected HasDifferences to be false")
	}
}

func TestSummarize_Counts(t *testing.T) {
	r := Result{
		Missing: []string{"A", "B"},
		Extra:   []string{"C"},
		Mismatched: []MismatchedEntry{
			{Key: "D", BaseValue: "x", OtherValue: "y"},
			{Key: "E", BaseValue: "1", OtherValue: "2"},
		},
	}
	s := Summarize(r)
	if s.MissingCount != 2 {
		t.Errorf("expected MissingCount 2, got %d", s.MissingCount)
	}
	if s.ExtraCount != 1 {
		t.Errorf("expected ExtraCount 1, got %d", s.ExtraCount)
	}
	if s.MismatchedCount != 2 {
		t.Errorf("expected MismatchedCount 2, got %d", s.MismatchedCount)
	}
	if s.TotalCount != 5 {
		t.Errorf("expected TotalCount 5, got %d", s.TotalCount)
	}
	if !s.HasDifferences() {
		t.Error("expected HasDifferences to be true")
	}
}

func TestSummary_String_NoDiff(t *testing.T) {
	s := Summary{}
	if s.String() != "no differences found" {
		t.Errorf("unexpected string: %q", s.String())
	}
}

func TestSummary_String_WithDiff(t *testing.T) {
	s := Summary{
		MissingCount:    1,
		ExtraCount:      2,
		MismatchedCount: 3,
		TotalCount:      6,
	}
	str := s.String()
	if !strings.Contains(str, "6 difference(s)") {
		t.Errorf("expected total count in string, got: %q", str)
	}
	if !strings.Contains(str, "1 missing") {
		t.Errorf("expected missing count in string, got: %q", str)
	}
	if !strings.Contains(str, "2 extra") {
		t.Errorf("expected extra count in string, got: %q", str)
	}
	if !strings.Contains(str, "3 mismatched") {
		t.Errorf("expected mismatched count in string, got: %q", str)
	}
}

func TestSummarize_OnlyMissing(t *testing.T) {
	r := Result{
		Missing: []string{"X"},
		Extra:   []string{},
		Mismatched: []MismatchedEntry{},
	}
	s := Summarize(r)
	if s.MissingCount != 1 {
		t.Errorf("expected MissingCount 1, got %d", s.MissingCount)
	}
	if s.ExtraCount != 0 {
		t.Errorf("expected ExtraCount 0, got %d", s.ExtraCount)
	}
	if s.MismatchedCount != 0 {
		t.Errorf("expected MismatchedCount 0, got %d", s.MismatchedCount)
	}
	if s.TotalCount != 1 {
		t.Errorf("expected TotalCount 1, got %d", s.TotalCount)
	}
	if !s.HasDifferences() {
		t.Error("expected HasDifferences to be true")
	}
}
