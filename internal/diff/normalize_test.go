package diff

import (
	"testing"
)

func TestNormalizeMap_TrimSpace(t *testing.T) {
	env := map[string]string{
		"HOST": "  localhost  ",
		"PORT": "8080",
	}
	opts := DefaultNormalizeOptions()
	out := NormalizeMap(env, opts)
	if out["HOST"] != "localhost" {
		t.Errorf("expected 'localhost', got %q", out["HOST"])
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected '8080', got %q", out["PORT"])
	}
}

func TestNormalizeMap_UppercaseKeys(t *testing.T) {
	env := map[string]string{
		"db_host": "localhost",
		"Api_Key": "secret",
	}
	opts := NormalizeOptions{UppercaseKeys: true}
	out := NormalizeMap(env, opts)
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST key after uppercase normalization")
	}
	if _, ok := out["API_KEY"]; !ok {
		t.Error("expected API_KEY key after uppercase normalization")
	}
	if _, ok := out["db_host"]; ok {
		t.Error("original lowercase key should not be present")
	}
}

func TestNormalizeMap_CollapseWhitespace(t *testing.T) {
	env := map[string]string{
		"MSG": "hello   world\there",
	}
	opts := NormalizeOptions{CollapseWhitespace: true}
	out := NormalizeMap(env, opts)
	if out["MSG"] != "hello world here" {
		t.Errorf("unexpected collapsed value: %q", out["MSG"])
	}
}

func TestNormalizeMap_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{
		"KEY": "  value  ",
	}
	opts := DefaultNormalizeOptions()
	_ = NormalizeMap(env, opts)
	if env["KEY"] != "  value  " {
		t.Error("NormalizeMap must not mutate the input map")
	}
}

func TestNormalizeMap_EmptyMap(t *testing.T) {
	out := NormalizeMap(map[string]string{}, DefaultNormalizeOptions())
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}

func TestCollapseWS_SingleSpaces(t *testing.T) {
	if got := collapseWS("a b c"); got != "a b c" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestCollapseWS_AllWhitespace(t *testing.T) {
	if got := collapseWS("   "); got != " " {
		t.Errorf("expected single space, got %q", got)
	}
}
