package diff

import "sort"

// EnvSet represents a set of environment variable keys.
type EnvSet struct {
	keys map[string]struct{}
}

// NewEnvSet creates an EnvSet from a map of env vars.
func NewEnvSet(env map[string]string) EnvSet {
	keys := make(map[string]struct{}, len(env))
	for k := range env {
		keys[k] = struct{}{}
	}
	return EnvSet{keys: keys}
}

// Has returns true if the key exists in the set.
func (s EnvSet) Has(key string) bool {
		_, ok := s.keys[key]
	return ok
}

// Keys returns a sorted slice of all keys in the set.
func (s EnvSet) Keys() []string {
	out := make([]string, 0, len(s.keys))
	for k := range s.keys {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Intersect returns keys present in both sets.
func (s EnvSet) Intersect(other EnvSet) []string {
	var result []string
	for k := range s.keys {
		if other.Has(k) {
			result = append(result, k)
		}
	}
	sort.Strings(result)
	return result
}

// Difference returns keys in s that are not in other.
func (s EnvSet) Difference(other EnvSet) []string {
	var result []string
	for k := range s.keys {
		if !other.Has(k) {
			result = append(result, k)
		}
	}
	sort.Strings(result)
	return result
}

// Len returns the number of keys in the set.
func (s EnvSet) Len() int {
	return len(s.keys)
}
