package diff

import (
	"testing"
)

func TestLint_NoIssues(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	result := Lint(env, DefaultLintOptions())
	if !result.Valid {
		t.Errorf("expected valid, got issues: %v", result.Issues)
	}
	if len(result.Issues) != 0 {
		t.Errorf("expected 0 issues, got %d", len(result.Issues))
	}
}

func TestLint_LowercaseKeyWarning(t *testing.T) {
	env := map[string]string{
		"app_name": "myapp",
	}
	opts := DefaultLintOptions()
	result := Lint(env, opts)
	if result.Valid {
		t.Error("expected invalid due to lowercase key")
	}
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != LintWarning {
		t.Errorf("expected warning, got %s", result.Issues[0].Severity)
	}
}

func TestLint_SpacesInValueError(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "  localhost",
	}
	result := Lint(env, DefaultLintOptions())
	if result.Valid {
		t.Error("expected invalid due to leading spaces in value")
	}
	found := false
	for _, issue := range result.Issues {
		if issue.Key == "DB_HOST" && issue.Severity == LintError {
			found = true
		}
	}
	if !found {
		t.Error("expected LintError for DB_HOST with leading spaces")
	}
}

func TestLint_DisabledRules(t *testing.T) {
	env := map[string]string{
		"lower_key": "  value  ",
	}
	opts := LintOptions{
		DisallowUppercaseOnly:  false,
		DisallowSpacesInValue: false,
		DisallowEmptyKeys:     false,
	}
	result := Lint(env, opts)
	if !result.Valid {
		t.Errorf("expected valid with all rules disabled, got: %v", result.Issues)
	}
}

func TestLint_MultipleIssuesSameKey(t *testing.T) {
	env := map[string]string{
		"bad key": "  oops  ",
	}
	opts := DefaultLintOptions()
	result := Lint(env, opts)
	// should have at least a warning (not uppercase) and an error (spaces in value)
	if len(result.Issues) < 2 {
		t.Errorf("expected at least 2 issues, got %d", len(result.Issues))
	}
}
