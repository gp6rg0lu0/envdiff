package diff

// RenameMap maps old key names to new key names.
type RenameMap map[string]string

// RenameResult holds the outcome of applying a rename map to an env map.
type RenameResult struct {
	// Renamed contains entries that were successfully renamed: old key -> new key.
	Renamed map[string]string
	// NotFound contains old keys from the rename map that were absent in the env.
	NotFound []string
	// Output is the env map after renames have been applied.
	Output map[string]string
}

// ApplyRenames takes an env map and a RenameMap, returns a new env map with
// keys renamed according to the map, plus a RenameResult describing what happened.
// Keys not present in the rename map are passed through unchanged.
// If an old key is not found in env, it is recorded in NotFound.
func ApplyRenames(env map[string]string, renames RenameMap) RenameResult {
	output := make(map[string]string, len(env))
	renamed := make(map[string]string)
	var notFound []string

	// Copy all keys first.
	for k, v := range env {
		output[k] = v
	}

	for oldKey, newKey := range renames {
		val, ok := output[oldKey]
		if !ok {
			notFound = append(notFound, oldKey)
			continue
		}
		delete(output, oldKey)
		output[newKey] = val
		renamed[oldKey] = newKey
	}

	return RenameResult{
		Renamed:  renamed,
		NotFound: notFound,
		Output:   output,
	}
}
