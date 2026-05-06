package diff

// Result holds the comparison result between two env files.
type Result struct {
	Missing   []string          // keys in base but not in target
	Extra     []string          // keys in target but not in base
	Mismatched map[string][2]string // key -> [baseVal, targetVal]
}

// Compare compares two parsed env maps (base vs target).
// base is the reference environment (e.g. .env.example).
// target is the environment being checked (e.g. .env.production).
func Compare(base, target map[string]string) Result {
	result := Result{
		Mismatched: make(map[string][2]string),
	}

	for key, baseVal := range base {
		targetVal, ok := target[key]
		if !ok {
			result.Missing = append(result.Missing, key)
		} else if baseVal != targetVal {
			result.Mismatched[key] = [2]string{baseVal, targetVal}
		}
	}

	for key := range target {
		if _, ok := base[key]; !ok {
			result.Extra = append(result.Extra, key)
		}
	}

	return result
}

// HasDifferences returns true if the result contains any differences.
func (r Result) HasDifferences() bool {
	return len(r.Missing) > 0 || len(r.Extra) > 0 || len(r.Mismatched) > 0
}
