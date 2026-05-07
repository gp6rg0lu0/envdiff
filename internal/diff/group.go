package diff

import "sort"

// GroupedResult holds diff entries organized by their status.
type GroupedResult struct {
	Missing    []Entry
	Extra      []Entry
	Mismatched []Entry
}

// Group partitions a Result into a GroupedResult by status.
func Group(r Result) GroupedResult {
	var g GroupedResult
	for _, e := range r.Entries {
		switch e.Status {
		case StatusMissing:
			g.Missing = append(g.Missing, e)
		case StatusExtra:
			g.Extra = append(g.Extra, e)
		case StatusMismatch:
			g.Mismatched = append(g.Mismatched, e)
		}
	}
	sortEntries(g.Missing)
	sortEntries(g.Extra)
	sortEntries(g.Mismatched)
	return g
}

// Keys returns all unique keys present across all groups, sorted.
func (g GroupedResult) Keys() []string {
	seen := make(map[string]struct{})
	for _, e := range g.Missing {
		seen[e.Key] = struct{}{}
	}
	for _, e := range g.Extra {
		seen[e.Key] = struct{}{}
	}
	for _, e := range g.Mismatched {
		seen[e.Key] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Total returns the total number of differing entries across all groups.
func (g GroupedResult) Total() int {
	return len(g.Missing) + len(g.Extra) + len(g.Mismatched)
}

func sortEntries(entries []Entry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
}
