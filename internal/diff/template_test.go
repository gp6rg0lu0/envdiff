package diff

import (
	"testing"
)

func TestExpandTemplates_NoRefs(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
	}
	res := ExpandTemplates(env)
	if len(res.Vars) != 0 {
		t.Errorf("expected no vars, got %d", len(res.Vars))
	}
	if len(res.Unresolved) != 0 {
		t.Errorf("expected no unresolved, got %d", len(res.Unresolved))
	}
}

func TestExpandTemplates_ResolvesKnownRef(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "http://${HOST}:${PORT}",
		"HOST":     "localhost",
		"PORT":     "9000",
	}
	res := ExpandTemplates(env)
	if len(res.Unresolved) != 0 {
		t.Errorf("expected all resolved, got unresolved: %v", res.Unresolved)
	}
	resolved := map[string]string{}
	for _, v := range res.Vars {
		resolved[v.Key] = v.Value
	}
	if resolved["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %s", resolved["HOST"])
	}
	if resolved["PORT"] != "9000" {
		t.Errorf("expected PORT=9000, got %s", resolved["PORT"])
	}
}

func TestExpandTemplates_UnresolvedRef(t *testing.T) {
	env := map[string]string{
		"DSN": "postgres://${DB_USER}:${DB_PASS}@localhost/db",
	}
	res := ExpandTemplates(env)
	if len(res.Unresolved) != 2 {
		t.Errorf("expected 2 unresolved, got %d: %v", len(res.Unresolved), res.Unresolved)
	}
}

func TestExpandTemplates_DeduplicatesRefs(t *testing.T) {
	env := map[string]string{
		"A": "${SHARED}",
		"B": "${SHARED}/extra",
	}
	res := ExpandTemplates(env)
	count := 0
	for _, v := range res.Vars {
		if v.Key == "SHARED" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected SHARED to appear once, got %d", count)
	}
}

func TestExpandTemplates_PartiallyResolved(t *testing.T) {
	env := map[string]string{
		"URL":  "${PROTO}://${MISSING_HOST}",
		"PROTO": "https",
	}
	res := ExpandTemplates(env)
	if len(res.Unresolved) != 1 || res.Unresolved[0] != "MISSING_HOST" {
		t.Errorf("expected MISSING_HOST unresolved, got %v", res.Unresolved)
	}
	var protoVar *TemplateVar
	for i := range res.Vars {
		if res.Vars[i].Key == "PROTO" {
			protoVar = &res.Vars[i]
		}
	}
	if protoVar == nil || !protoVar.Resolved || protoVar.Value != "https" {
		t.Errorf("expected PROTO resolved to https")
	}
}
