package diff

import "fmt"

// SchemaRule defines an expected key with optional type and pattern constraints.
type SchemaRule struct {
	Key      string
	Required bool
	Type     string // "string", "int", "bool", "url"
	Pattern  string // optional regex hint (stored for reporting)
}

// SchemaResult holds the outcome of a schema validation pass.
type SchemaResult struct {
	Valid    bool
	Violations []SchemaViolation
}

// SchemaViolation describes a single schema rule failure.
type SchemaViolation struct {
	Key     string
	Rule    string
	Message string
}

// ValidateSchema checks the provided env map against a set of SchemaRules.
func ValidateSchema(env map[string]string, rules []SchemaRule) SchemaResult {
	var violations []SchemaViolation

	for _, rule := range rules {
		val, exists := env[rule.Key]

		if rule.Required && !exists {
			violations = append(violations, SchemaViolation{
				Key:     rule.Key,
				Rule:    "required",
				Message: fmt.Sprintf("key %q is required but missing", rule.Key),
			})
			continue
		}

		if !exists {
			continue
		}

		if rule.Type != "" {
			if err := checkType(val, rule.Type); err != nil {
				violations = append(violations, SchemaViolation{
					Key:     rule.Key,
					Rule:    "type",
					Message: fmt.Sprintf("key %q: %s", rule.Key, err.Error()),
				})
			}
		}
	}

	return SchemaResult{
		Valid:      len(violations) == 0,
		Violations: violations,
	}
}

func checkType(val, typ string) error {
	switch typ {
	case "int":
		var n int
		_, err := fmt.Sscanf(val, "%d", &n)
		if err != nil || fmt.Sprintf("%d", n) != val {
			return fmt.Errorf("expected integer, got %q", val)
		}
	case "bool":
		switch val {
		case "true", "false", "1", "0", "yes", "no":
		default:
			return fmt.Errorf("expected boolean, got %q", val)
		}
	case "url":
		if len(val) < 7 || (val[:7] != "http://" && (len(val) < 8 || val[:8] != "https://")) {
			return fmt.Errorf("expected URL starting with http:// or https://, got %q", val)
		}
	}
	return nil
}
