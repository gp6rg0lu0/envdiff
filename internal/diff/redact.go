package diff

import "strings"

// redactedPlaceholder is used to replace sensitive values.
const redactedPlaceholder = "[REDACTED]"

// sensitivePatterns are substrings that indicate a key holds sensitive data.
var sensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY",
	"PRIVATE", "CREDENTIAL", "AUTH", "CERT", "KEY",
}

// IsSensitive returns true if the key name matches any sensitive pattern.
func IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range sensitivePatterns {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

// RedactResult returns a copy of the Result with sensitive values replaced.
func RedactResult(r Result) Result {
	redacted := Result{
		Missing:    r.Missing,
		Extra:      r.Extra,
		Mismatched: make(map[string][2]string, len(r.Mismatched)),
	}

	for key, pair := range r.Mismatched {
		if IsSensitive(key) {
			redacted.Mismatched[key] = [2]string{redactedPlaceholder, redactedPlaceholder}
		} else {
			redacted.Mismatched[key] = pair
		}
	}

	return redacted
}
