package diff

import (
	"strings"
)

// NormalizeOptions controls how env values are normalized before comparison.
type NormalizeOptions struct {
	// TrimSpace removes leading/trailing whitespace from values.
	TrimSpace bool
	// LowercaseKeys normalizes all keys to uppercase before comparison.
	UppercaseKeys bool
	// CollapseWhitespace replaces internal runs of whitespace with a single space.
	CollapseWhitespace bool
}

// DefaultNormalizeOptions returns a sensible default configuration.
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		TrimSpace:          true,
		UppercaseKeys:      false,
		CollapseWhitespace: false,
	}
}

// NormalizeMap applies normalization rules to a copy of the provided env map.
func NormalizeMap(env map[string]string, opts NormalizeOptions) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		nk := k
		if opts.UppercaseKeys {
			nk = strings.ToUpper(k)
		}
		nv := v
		if opts.TrimSpace {
			nv = strings.TrimSpace(nv)
		}
		if opts.CollapseWhitespace {
			nv = collapseWS(nv)
		}
		out[nk] = nv
	}
	return out
}

// collapseWS replaces runs of whitespace characters with a single space.
func collapseWS(s string) string {
	var b strings.Builder
	inSpace := false
	for _, r := range s {
		if r == ' ' || r == '\t' {
			if !inSpace {
				b.WriteRune(' ')
				inSpace = true
			}
		} else {
			b.WriteRune(r)
			inSpace = false
		}
	}
	return b.String()
}
