package diff

// FilterOptions controls which difference types are included in a result.
type FilterOptions struct {
	ShowMissing   bool
	ShowExtra     bool
	ShowMismatched bool
}

// DefaultFilterOptions returns options that show all difference types.
func DefaultFilterOptions() FilterOptions {
	return FilterOptions{
		ShowMissing:    true,
		ShowExtra:      true,
		ShowMismatched: true,
	}
}

// Filter returns a new Result containing only the entries matching the given options.
func Filter(r Result, opts FilterOptions) Result {
	out := Result{
		Missing:    []string{},
		Extra:      []string{},
		Mismatched: []MismatchedEntry{},
	}

	if opts.ShowMissing {
		out.Missing = append(out.Missing, r.Missing...)
	}

	if opts.ShowExtra {
		out.Extra = append(out.Extra, r.Extra...)
	}

	if opts.ShowMismatched {
		out.Mismatched = append(out.Mismatched, r.Mismatched...)
	}

	return out
}

// IsEmpty returns true when the result contains no differences.
func IsEmpty(r Result) bool {
	return len(r.Missing) == 0 && len(r.Extra) == 0 && len(r.Mismatched) == 0
}
