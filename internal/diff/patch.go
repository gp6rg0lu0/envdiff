package diff

import (
	"fmt"
	"sort"
	"strings"
)

// PatchEntry represents a single line in a patch/diff output.
type PatchEntry struct {
	Key    string
	OldVal string
	NewVal string
	Op     string // "add", "remove", "change"
}

// Patch generates a list of patch entries describing how to transform
// the base environment into the target environment based on a Result.
func Patch(result Result) []PatchEntry {
	var entries []PatchEntry

	for _, key := range sortedKeys(result.Missing) {
		entries = append(entries, PatchEntry{
			Key:    key,
			NewVal: result.Missing[key],
			Op:     "add",
		})
	}

	for _, key := range sortedKeys(result.Extra) {
		entries = append(entries, PatchEntry{
			Key:    key,
			OldVal: result.Extra[key],
			Op:     "remove",
		})
	}

	for _, key := range sortedKeys(result.Mismatched) {
		pair := result.Mismatched[key]
		entries = append(entries, PatchEntry{
			Key:    key,
			OldVal: pair[0],
			NewVal: pair[1],
			Op:     "change",
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return entries
}

// FormatPatch renders patch entries as a unified-diff-style string.
func FormatPatch(entries []PatchEntry) string {
	var sb strings.Builder
	for _, e := range entries {
		switch e.Op {
		case "add":
			sb.WriteString(fmt.Sprintf("+ %s=%s\n", e.Key, e.NewVal))
		case "remove":
			sb.WriteString(fmt.Sprintf("- %s=%s\n", e.Key, e.OldVal))
		case "change":
			sb.WriteString(fmt.Sprintf("~ %s: %s -> %s\n", e.Key, e.OldVal, e.NewVal))
		}
	}
	return sb.String()
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
