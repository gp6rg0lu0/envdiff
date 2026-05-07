package diff

import "fmt"

// Summary holds aggregate counts of differences in a Result.
type Summary struct {
	MissingCount    int
	ExtraCount      int
	MismatchedCount int
	TotalCount      int
}

// Summarize computes a Summary from a Result.
func Summarize(r Result) Summary {
	s := Summary{
		MissingCount:    len(r.Missing),
		ExtraCount:      len(r.Extra),
		MismatchedCount: len(r.Mismatched),
	}
	s.TotalCount = s.MissingCount + s.ExtraCount + s.MismatchedCount
	return s
}

// String returns a human-readable one-line summary.
func (s Summary) String() string {
	if s.TotalCount == 0 {
		return "no differences found"
	}
	return fmt.Sprintf(
		"%d difference(s): %d missing, %d extra, %d mismatched",
		s.TotalCount, s.MissingCount, s.ExtraCount, s.MismatchedCount,
	)
}

// HasDifferences returns true when TotalCount is greater than zero.
func (s Summary) HasDifferences() bool {
	return s.TotalCount > 0
}
