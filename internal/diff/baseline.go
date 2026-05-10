package diff

import "sort"

// BaselineEntry represents a single key captured in a baseline snapshot.
type BaselineEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Baseline is a named snapshot of an environment's key-value pairs.
type Baseline struct {
	Name    string          `json:"name"`
	Entries []BaselineEntry `json:"entries"`
}

// NewBaseline creates a Baseline from a parsed env map.
func NewBaseline(name string, env map[string]string) Baseline {
	entries := make([]BaselineEntry, 0, len(env))
	for k, v := range env {
		entries = append(entries, BaselineEntry{Key: k, Value: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return Baseline{Name: name, Entries: entries}
}

// ToMap converts a Baseline back into a key-value map.
func (b Baseline) ToMap() map[string]string {
	m := make(map[string]string, len(b.Entries))
	for _, e := range b.Entries {
		m[e.Key] = e.Value
	}
	return m
}

// DiffBaseline compares the current env map against a saved Baseline and
// returns a Result describing what has changed since the baseline was taken.
func DiffBaseline(baseline Baseline, current map[string]string) Result {
	return Compare(baseline.ToMap(), current)
}
