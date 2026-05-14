package diff

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func makeEnvSet(keys ...string) EnvSet {
	m := make(map[string]string, len(keys))
	for _, k := range keys {
		m[k] = "value"
	}
	return NewEnvSet(m)
}

func TestEnvSet_Has(t *testing.T) {
	s := makeEnvSet("FOO", "BAR")
	if !s.Has("FOO") {
		t.Error("expected FOO to be present")
	}
	if s.Has("MISSING") {
		t.Error("expected MISSING to be absent")
	}
}

func TestEnvSet_Keys_Sorted(t *testing.T) {
	s := makeEnvSet("ZEBRA", "ALPHA", "MIDDLE")
	got := s.Keys()
	want := []string{"ALPHA", "MIDDLE", "ZEBRA"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Keys() mismatch (-want +got):\n%s", diff)
	}
}

func TestEnvSet_Len(t *testing.T) {
	s := makeEnvSet("A", "B", "C")
	if s.Len() != 3 {
		t.Errorf("expected Len 3, got %d", s.Len())
	}
}

func TestEnvSet_Intersect(t *testing.T) {
	a := makeEnvSet("FOO", "BAR", "SHARED")
	b := makeEnvSet("SHARED", "EXTRA", "BAR")
	got := a.Intersect(b)
	want := []string{"BAR", "SHARED"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Intersect() mismatch (-want +got):\n%s", diff)
	}
}

func TestEnvSet_Difference(t *testing.T) {
	a := makeEnvSet("FOO", "BAR", "ONLY_A")
	b := makeEnvSet("FOO", "BAR", "ONLY_B")
	got := a.Difference(b)
	want := []string{"ONLY_A"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Difference() mismatch (-want +got):\n%s", diff)
	}
}

func TestEnvSet_Intersect_Empty(t *testing.T) {
	a := makeEnvSet("FOO")
	b := makeEnvSet("BAR")
	got := a.Intersect(b)
	if len(got) != 0 {
		t.Errorf("expected empty intersection, got %v", got)
	}
}

func TestEnvSet_Difference_AllShared(t *testing.T) {
	a := makeEnvSet("FOO", "BAR")
	b := makeEnvSet("FOO", "BAR")
	got := a.Difference(b)
	if len(got) != 0 {
		t.Errorf("expected empty difference, got %v", got)
	}
}
