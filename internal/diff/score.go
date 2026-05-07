package diff

import (
	"fmt"

	"github.com/user/envdiff/internal/parser"
)

// Score represents a compatibility score between two env files.
type Score struct {
	Total      int
	Matched    int
	Missing    int
	Extra      int
	Mismatched int
	Percent    float64
}

// String returns a human-readable representation of the score.
func (s Score) String() string {
	return fmt.Sprintf("%.1f%% compatible (%d/%d keys matched, %d missing, %d extra, %d mismatched)",
		s.Percent, s.Matched, s.Total, s.Missing, s.Extra, s.Mismatched)
}

// ScoreEnvs computes a compatibility score between a base and target env map.
// The score is based on how many keys from the union of both files match.
func ScoreEnvs(base, target parser.Env) Score {
	result := Compare(base, target)

	unionKeys := make(map[string]struct{})
	for k := range base {
		unionKeys[k] = struct{}{}
	}
	for k := range target {
		unionKeys[k] = struct{}{}
	}

	total := len(unionKeys)
	if total == 0 {
		return Score{Percent: 100.0}
	}

	missing := len(result.Missing)
	extra := len(result.Extra)
	mismatched := len(result.Mismatched)
	matched := total - missing - extra - mismatched
	if matched < 0 {
		matched = 0
	}

	percent := float64(matched) / float64(total) * 100.0

	return Score{
		Total:      total,
		Matched:    matched,
		Missing:    missing,
		Extra:      extra,
		Mismatched: mismatched,
		Percent:    percent,
	}
}
