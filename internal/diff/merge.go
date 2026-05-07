package diff

// MergeResult holds the outcome of merging a base env map with overrides.
type MergeResult struct {
	// Merged is the final key-value map after applying overrides.
	Merged map[string]string
	// Added contains keys that were added from the override (not in base).
	Added []string
	// Overridden contains keys whose values were changed by the override.
	Overridden map[string]Override
}

// Override captures the before/after values for a key that was changed.
type Override struct {
	Base     string
	Override string
}

// Merge combines a base env map with one or more override maps, applying them
// in order. Later overrides take precedence over earlier ones.
func Merge(base map[string]string, overrides ...map[string]string) MergeResult {
	merged := make(map[string]string, len(base))
	for k, v := range base {
		merged[k] = v
	}

	added := []string{}
	overridden := map[string]Override{}

	for _, ov := range overrides {
		for k, v := range ov {
			baseVal, exists := merged[k]
			if !exists {
				added = append(added, k)
				merged[k] = v
			} else if baseVal != v {
				overridden[k] = Override{Base: baseVal, Override: v}
				merged[k] = v
			}
		}
	}

	return MergeResult{
		Merged:     merged,
		Added:      added,
		Overridden: overridden,
	}
}
