package diff

import "strings"

// LintSeverity represents the severity level of a lint issue.
type LintSeverity string

const (
	LintError   LintSeverity = "error"
	LintWarning LintSeverity = "warning"
	LintInfo    LintSeverity = "info"
)

// LintIssue describes a single linting problem found in an env map.
type LintIssue struct {
	Key      string
	Message  string
	Severity LintSeverity
}

// LintResult holds all issues found during linting.
type LintResult struct {
	Issues []LintIssue
	Valid  bool
}

// LintOptions controls which lint rules are applied.
type LintOptions struct {
	DisallowUppercaseOnly bool // warn if key is not UPPER_SNAKE_CASE
	DisallowSpacesInValue bool // error if value contains leading/trailing spaces
	DisallowEmptyKeys     bool // error if key is empty string
	DisallowDuplicates    bool // error on duplicate keys (parser-level; best-effort)
}

// DefaultLintOptions returns a sensible default configuration.
func DefaultLintOptions() LintOptions {
	return LintOptions{
		DisallowUppercaseOnly:  true,
		DisallowSpacesInValue: true,
		DisallowEmptyKeys:     true,
	}
}

// Lint inspects an env map and returns a LintResult with any issues found.
func Lint(env map[string]string, opts LintOptions) LintResult {
	var issues []LintIssue

	for k, v := range env {
		if opts.DisallowEmptyKeys && strings.TrimSpace(k) == "" {
			issues = append(issues, LintIssue{
				Key:      k,
				Message:  "key must not be empty or whitespace",
				Severity: LintError,
			})
		}

		if opts.DisallowUppercaseOnly && k != strings.ToUpper(k) {
			issues = append(issues, LintIssue{
				Key:      k,
				Message:  "key should be UPPER_SNAKE_CASE",
				Severity: LintWarning,
			})
		}

		if opts.DisallowSpacesInValue && v != strings.TrimSpace(v) {
			issues = append(issues, LintIssue{
				Key:      k,
				Message:  "value has leading or trailing whitespace",
				Severity: LintError,
			})
		}
	}

	return LintResult{
		Issues: issues,
		Valid:  len(issues) == 0,
	}
}
