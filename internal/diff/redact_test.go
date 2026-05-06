package diff_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/diff"
)

func TestIsSensitive(t *testing.T) {
	cases := []struct {
		key      string
		want     bool
	}{
		{"DB_PASSWORD", true},
		{"API_KEY", true},
		{"SECRET_TOKEN", true},
		{"AUTH_HEADER", true},
		{"DATABASE_URL", false},
		{"PORT", false},
		{"APP_ENV", false},
		{"PRIVATE_KEY", true},
	}

	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			got := diff.IsSensitive(tc.key)
			if got != tc.want {
				t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.want)
			}
		})
	}
}

func TestRedactResult_RedactsSensitiveKeys(t *testing.T) {
	r := diff.Result{
		Mismatched: map[string][2]string{
			"DB_PASSWORD": {"secret1", "secret2"},
			"APP_NAME":    {"myapp", "otherapp"},
		},
	}

	redacted := diff.RedactResult(r)

	pair, ok := redacted.Mismatched["DB_PASSWORD"]
	if !ok {
		t.Fatal("expected DB_PASSWORD in mismatched after redaction")
	}
	if pair[0] != "[REDACTED]" || pair[1] != "[REDACTED]" {
		t.Errorf("expected redacted values, got %v", pair)
	}

	appPair := redacted.Mismatched["APP_NAME"]
	if appPair[0] != "myapp" || appPair[1] != "otherapp" {
		t.Errorf("non-sensitive key should not be redacted, got %v", appPair)
	}
}

func TestRedactResult_PreservesMissingAndExtra(t *testing.T) {
	r := diff.Result{
		Missing:    []string{"MISSING_KEY"},
		Extra:      []string{"EXTRA_KEY"},
		Mismatched: map[string][2]string{},
	}

	redacted := diff.RedactResult(r)
	if len(redacted.Missing) != 1 || redacted.Missing[0] != "MISSING_KEY" {
		t.Errorf("Missing keys should be preserved: %v", redacted.Missing)
	}
	if len(redacted.Extra) != 1 || redacted.Extra[0] != "EXTRA_KEY" {
		t.Errorf("Extra keys should be preserved: %v", redacted.Extra)
	}
}
