package diff

import (
	"testing"
)

func TestValidate_NoIssues(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	result := Validate(env, DefaultValidateOptions())
	if !result.Valid() {
		t.Errorf("expected no issues, got %v", result.Issues)
	}
}

func TestValidate_EmptyValueDisallowed(t *testing.T) {
	env := map[string]string{"HOST": "", "PORT": "8080"}
	opts := ValidateOptions{DisallowEmptyValues: true}
	result := Validate(env, opts)
	if result.Valid() {
		t.Fatal("expected validation to fail for empty value")
	}
	if len(result.Issues) != 1 || result.Issues[0].Key != "HOST" {
		t.Errorf("unexpected issues: %v", result.Issues)
	}
}

func TestValidate_RequiredKeyMissing(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	opts := ValidateOptions{RequiredKeys: []string{"PORT", "HOST"}}
	result := Validate(env, opts)
	if result.Valid() {
		t.Fatal("expected validation to fail for missing required key")
	}
	if len(result.Issues) != 1 || result.Issues[0].Key != "PORT" {
		t.Errorf("unexpected issues: %v", result.Issues)
	}
}

func TestValidate_ForbiddenKeyPresent(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "DEBUG": "true"}
	opts := ValidateOptions{ForbiddenKeys: []string{"DEBUG"}}
	result := Validate(env, opts)
	if result.Valid() {
		t.Fatal("expected validation to fail for forbidden key")
	}
	if len(result.Issues) != 1 || result.Issues[0].Key != "DEBUG" {
		t.Errorf("unexpected issues: %v", result.Issues)
	}
}

func TestValidate_MultipleIssues(t *testing.T) {
	env := map[string]string{"HOST": "", "SECRET": "abc"}
	opts := ValidateOptions{
		DisallowEmptyValues: true,
		RequiredKeys:        []string{"PORT"},
		ForbiddenKeys:       []string{"SECRET"},
	}
	result := Validate(env, opts)
	if len(result.Issues) != 3 {
		t.Errorf("expected 3 issues, got %d: %v", len(result.Issues), result.Issues)
	}
}

func TestValidationIssue_String(t *testing.T) {
	v := ValidationIssue{Key: "FOO", Message: "value is empty"}
	if v.String() != "FOO: value is empty" {
		t.Errorf("unexpected string: %s", v.String())
	}
}
