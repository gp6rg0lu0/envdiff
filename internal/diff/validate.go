package diff

import (
	"fmt"
	"strings"
)

// ValidationIssue represents a single validation problem found in an env map.
type ValidationIssue struct {
	Key     string
	Message string
}

func (v ValidationIssue) String() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// ValidationResult holds all issues found during validation.
type ValidationResult struct {
	Issues []ValidationIssue
}

func (r ValidationResult) Valid() bool {
	return len(r.Issues) == 0
}

// ValidateOptions controls which checks are performed.
type ValidateOptions struct {
	DisallowEmptyValues bool
	RequiredKeys        []string
	ForbiddenKeys       []string
}

// DefaultValidateOptions returns sensible defaults.
func DefaultValidateOptions() ValidateOptions {
	return ValidateOptions{
		DisallowEmptyValues: false,
		RequiredKeys:        nil,
		ForbiddenKeys:       nil,
	}
}

// Validate checks an env map for common issues based on the provided options.
func Validate(env map[string]string, opts ValidateOptions) ValidationResult {
	var issues []ValidationIssue

	if opts.DisallowEmptyValues {
		for k, v := range env {
			if strings.TrimSpace(v) == "" {
				issues = append(issues, ValidationIssue{Key: k, Message: "value is empty"})
			}
		}
	}

	for _, req := range opts.RequiredKeys {
		if _, ok := env[req]; !ok {
			issues = append(issues, ValidationIssue{Key: req, Message: "required key is missing"})
		}
	}

	for _, forbidden := range opts.ForbiddenKeys {
		if _, ok := env[forbidden]; ok {
			issues = append(issues, ValidationIssue{Key: forbidden, Message: "forbidden key is present"})
		}
	}

	return ValidationResult{Issues: issues}
}
