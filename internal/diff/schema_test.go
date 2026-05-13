package diff

import "testing"

func TestValidateSchema_NoViolations(t *testing.T) {
	env := map[string]string{
		"PORT":     "8080",
		"DEBUG":    "true",
		"BASE_URL": "https://example.com",
	}
	rules := []SchemaRule{
		{Key: "PORT", Required: true, Type: "int"},
		{Key: "DEBUG", Required: true, Type: "bool"},
		{Key: "BASE_URL", Required: false, Type: "url"},
	}
	res := ValidateSchema(env, rules)
	if !res.Valid {
		t.Fatalf("expected valid, got violations: %v", res.Violations)
	}
}

func TestValidateSchema_MissingRequiredKey(t *testing.T) {
	env := map[string]string{}
	rules := []SchemaRule{
		{Key: "DATABASE_URL", Required: true},
	}
	res := ValidateSchema(env, rules)
	if res.Valid {
		t.Fatal("expected invalid")
	}
	if len(res.Violations) != 1 || res.Violations[0].Rule != "required" {
		t.Fatalf("unexpected violations: %v", res.Violations)
	}
}

func TestValidateSchema_TypeMismatch_Int(t *testing.T) {
	env := map[string]string{"PORT": "not-a-number"}
	rules := []SchemaRule{{Key: "PORT", Required: true, Type: "int"}}
	res := ValidateSchema(env, rules)
	if res.Valid {
		t.Fatal("expected invalid")
	}
	if res.Violations[0].Rule != "type" {
		t.Fatalf("expected type violation, got %q", res.Violations[0].Rule)
	}
}

func TestValidateSchema_TypeMismatch_Bool(t *testing.T) {
	env := map[string]string{"FLAG": "maybe"}
	rules := []SchemaRule{{Key: "FLAG", Type: "bool"}}
	res := ValidateSchema(env, rules)
	if res.Valid {
		t.Fatal("expected invalid")
	}
}

func TestValidateSchema_TypeMismatch_URL(t *testing.T) {
	env := map[string]string{"ENDPOINT": "ftp://bad.com"}
	rules := []SchemaRule{{Key: "ENDPOINT", Type: "url"}}
	res := ValidateSchema(env, rules)
	if res.Valid {
		t.Fatal("expected invalid")
	}
}

func TestValidateSchema_OptionalMissingKeySkipped(t *testing.T) {
	env := map[string]string{}
	rules := []SchemaRule{{Key: "OPTIONAL_KEY", Required: false, Type: "int"}}
	res := ValidateSchema(env, rules)
	if !res.Valid {
		t.Fatalf("expected valid for optional missing key, got: %v", res.Violations)
	}
}

func TestValidateSchema_MultipleViolations(t *testing.T) {
	env := map[string]string{"PORT": "abc"}
	rules := []SchemaRule{
		{Key: "PORT", Required: true, Type: "int"},
		{Key: "SECRET", Required: true},
	}
	res := ValidateSchema(env, rules)
	if res.Valid {
		t.Fatal("expected invalid")
	}
	if len(res.Violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(res.Violations))
	}
}
