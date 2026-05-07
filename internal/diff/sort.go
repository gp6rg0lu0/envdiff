package diff

import (
	"sort"

	"github.com/user/envdiff/internal/parser"
)

// SortOrder defines how diff results should be sorted.
type SortOrder string

const (
	SortByKey      SortOrder = "key"
	SortByStatus   SortOrder = "status"
	SortByStatusKey SortOrder = "status-key"
)

// statusRank returns a numeric rank for a given diff status for ordering.
func statusRank(d parser.DiffEntry) int {
	switch {
	case d.Missing:
		return 0
	case d.Extra:
		return 1
	case d.Mismatched:
		return 2
	default:
		return 3
	}
}

// SortResult returns a new Result with Differences sorted by the given order.
// The original Result is not mutated.
func SortResult(r parser.Result, order SortOrder) parser.Result {
	entries := make([]parser.DiffEntry, len(r.Differences))
	copy(entries, r.Differences)

	switch order {
	case SortByKey:
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Key < entries[j].Key
		})
	case SortByStatus:
		sort.SliceStable(entries, func(i, j int) bool {
			return statusRank(entries[i]) < statusRank(entries[j])
		})
	case SortByStatusKey:
		sort.Slice(entries, func(i, j int) bool {
			ri, rj := statusRank(entries[i]), statusRank(entries[j])
			if ri != rj {
				return ri < rj
			}
			return entries[i].Key < entries[j].Key
		})
	default:
		// no-op: return as-is
	}

	return parser.Result{
		BaseFile:    r.BaseFile,
		CompareFile: r.CompareFile,
		Differences: entries,
	}
}
