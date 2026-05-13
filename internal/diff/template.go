package diff

import (
	"fmt"
	"strings"
)

// TemplateVar represents a variable reference found in a .env value.
type TemplateVar struct {
	Key      string
	Raw      string
	Resolved bool
	Value    string
}

// TemplateResult holds the analysis of variable interpolation in an env map.
type TemplateResult struct {
	Vars     []TemplateVar
	Unresolved []string
}

// ExpandTemplates scans env values for ${VAR} or $VAR references and
// attempts to resolve them using the provided env map.
func ExpandTemplates(env map[string]string) TemplateResult {
	var result TemplateResult
	seen := map[string]bool{}

	for _, val := range env {
		refs := extractRefs(val)
		for _, ref := range refs {
			if seen[ref] {
				continue
			}
			seen[ref] = true

			tv := TemplateVar{
				Key: ref,
				Raw: fmt.Sprintf("${%s}", ref),
			}
			if resolved, ok := env[ref]; ok {
				tv.Resolved = true
				tv.Value = resolved
			} else {
				result.Unresolved = append(result.Unresolved, ref)
			}
			result.Vars = append(result.Vars, tv)
		}
	}
	return result
}

// extractRefs finds all ${VAR} style references in a string.
func extractRefs(s string) []string {
	var refs []string
	for {
		start := strings.Index(s, "${")
		if start == -1 {
			break
		}
		end := strings.Index(s[start:], "}")
		if end == -1 {
			break
		}
		key := s[start+2 : start+end]
		if key != "" {
			refs = append(refs, key)
		}
		s = s[start+end+1:]
	}
	return refs
}
