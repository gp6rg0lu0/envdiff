package diff

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
)

func TestScoreEnvs_PerfectMatch(t *testing.T) {
	base := parser.Env{"A": "1", "B": "2"}
	target := parser.Env{"A": "1", "B": "2"}

	s := ScoreEnvs(base, target)
	if s.Percent != 100.0 {
		t.Errorf("expected 100%%, got %.1f%%", s.Percent)
	}
	if s.Matched != 2 {
		t.Errorf("expected 2 matched, got %d", s.Matched)
	}
}

func TestScoreEnvs_Empty(t *testing.T) {
	s := ScoreEnvs(parser.Env{}, parser.Env{})
	if s.Percent != 100.0 {
		t.Errorf("expected 100%% for empty envs, got %.1f%%", s.Percent)
	}
}

func TestScoreEnvs_AllMissing(t *testing.T) {
	base := parser.Env{"A": "1", "B": "2"}
	target := parser.Env{}

	s := ScoreEnvs(base, target)
	if s.Missing != 2 {
		t.Errorf("expected 2 missing, got %d", s.Missing)
	}
	if s.Percent != 0.0 {
		t.Errorf("expected 0%%, got %.1f%%", s.Percent)
	}
}

func TestScoreEnvs_MixedDifferences(t *testing.T) {
	base := parser.Env{"A": "1", "B": "2", "C": "3"}
	target := parser.Env{"A": "1", "B": "changed", "D": "4"}

	// union: A, B, C, D = 4 keys
	// matched: A = 1
	// missing: C = 1
	// extra: D = 1
	// mismatched: B = 1
	s := ScoreEnvs(base, target)
	if s.Total != 4 {
		t.Errorf("expected total 4, got %d", s.Total)
	}
	if s.Matched != 1 {
		t.Errorf("expected 1 matched, got %d", s.Matched)
	}
	if s.Missing != 1 {
		t.Errorf("expected 1 missing, got %d", s.Missing)
	}
	if s.Extra != 1 {
		t.Errorf("expected 1 extra, got %d", s.Extra)
	}
	if s.Mismatched != 1 {
		t.Errorf("expected 1 mismatched, got %d", s.Mismatched)
	}
	expectedPct := 25.0
	if s.Percent != expectedPct {
		t.Errorf("expected %.1f%%, got %.1f%%", expectedPct, s.Percent)
	}
}

func TestScoreEnvs_AllExtra(t *testing.T) {
	base := parser.Env{}
	target := parser.Env{"X": "1", "Y": "2"}

	s := ScoreEnvs(base, target)
	if s.Extra != 2 {
		t.Errorf("expected 2 extra, got %d", s.Extra)
	}
	if s.Matched != 0 {
		t.Errorf("expected 0 matched, got %d", s.Matched)
	}
	if s.Percent != 0.0 {
		t.Errorf("expected 0%%, got %.1f%%", s.Percent)
	}
}

func TestScore_String(t *testing.T) {
	s := Score{
		Total: 4, Matched: 2, Missing: 1, Extra: 0, Mismatched: 1, Percent: 50.0,
	}
	got := s.String()
	if got == "" {
		t.Error("expected non-empty string from Score.String()")
	}
}
