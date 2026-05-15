package diff

import (
	"strings"
)

// MaskOptions controls how values are masked in output.
type MaskOptions struct {
	// VisibleChars is the number of characters to reveal at the start of a value.
	VisibleChars int
	// MaskChar is the character used to replace hidden characters.
	MaskChar rune
	// MinLength is the minimum value length to partially reveal; shorter values are fully masked.
	MinLength int
}

// DefaultMaskOptions returns sensible defaults for masking.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		VisibleChars: 3,
		MaskChar:     '*',
		MinLength:    6,
	}
}

// MaskValue partially or fully masks a string value according to opts.
// If the value is shorter than MinLength, it is fully masked.
// Otherwise, VisibleChars characters are shown followed by mask characters.
func MaskValue(value string, opts MaskOptions) string {
	if len(value) == 0 {
		return ""
	}
	if len(value) < opts.MinLength {
		return strings.Repeat(string(opts.MaskChar), len(value))
	}
	visible := opts.VisibleChars
	if visible >= len(value) {
		visible = len(value) - 1
	}
	maskLen := len(value) - visible
	if maskLen < 1 {
		maskLen = 1
	}
	return value[:visible] + strings.Repeat(string(opts.MaskChar), maskLen)
}

// MaskResult returns a copy of result where sensitive values are masked
// using the provided options. Non-sensitive keys are left unchanged.
func MaskResult(result Result, opts MaskOptions) Result {
	masked := make([]Entry, len(result.Entries))
	for i, e := range result.Entries {
		if IsSensitive(e.Key) {
			if e.BaseValue != "" {
				e.BaseValue = MaskValue(e.BaseValue, opts)
			}
			if e.OtherValue != "" {
				e.OtherValue = MaskValue(e.OtherValue, opts)
			}
		}
		masked[i] = e
	}
	return Result{Entries: masked}
}
