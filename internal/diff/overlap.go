package diff

import "sort"

// OverlapResult holds the analysis of key overlap between two env maps.
type OverlapResult struct {
	// SharedKeys are keys present in both base and compare.
	SharedKeys []string
	// UniqueToBase are keys only in base.
	UniqueToBase []string
	// UniqueToCompare are keys only in compare.
	UniqueToCompare []string
	// OverlapPercent is the Jaccard similarity: shared / union * 100.
	OverlapPercent float64
}

// Overlap computes key overlap statistics between two env maps.
func Overlap(base, compare map[string]string) OverlapResult {
	shared := []string{}
	uniqueBase := []string{}
	uniqueCompare := []string{}

	compareSet := make(map[string]bool, len(compare))
	for k := range compare {
		compareSet[k] = true
	}

	for k := range base {
		if compareSet[k] {
			shared = append(shared, k)
		} else {
			uniqueBase = append(uniqueBase, k)
		}
	}

	for k := range compare {
		if _, ok := base[k]; !ok {
			uniqueCompare = append(uniqueCompare, k)
		}
	}

	sort.Strings(shared)
	sort.Strings(uniqueBase)
	sort.Strings(uniqueCompare)

	union := len(shared) + len(uniqueBase) + len(uniqueCompare)
	var pct float64
	if union > 0 {
		pct = float64(len(shared)) / float64(union) * 100.0
	}

	return OverlapResult{
		SharedKeys:      shared,
		UniqueToBase:    uniqueBase,
		UniqueToCompare: uniqueCompare,
		OverlapPercent:  pct,
	}
}
