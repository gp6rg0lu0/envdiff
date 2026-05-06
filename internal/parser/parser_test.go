package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseFile_Basic(t *testing.T) {
	path := writeTempEnv(t, "KEY1=value1\nKEY2=value2\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY1"] != "value1" || env["KEY2"] != "value2" {
		t.Errorf("unexpected env: %v", env)
	}
}

func TestParseFile_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# comment\n\nKEY=val\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 || env["KEY"] != "val" {
		t.Errorf("expected only KEY=val, got: %v", env)
	}
}

func TestParseFile_StripQuotes(t *testing.T) {
	path := writeTempEnv(t, `DOUBLE="hello world"
SINGLE='goodbye'
`)
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["DOUBLE"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", env["DOUBLE"])
	}
	if env["SINGLE"] != "goodbye" {
		t.Errorf("expected 'goodbye', got %q", env["SINGLE"])
	}
}

func TestParseFile_MalformedLine(t *testing.T) {
	path := writeTempEnv(t, "NODIVIDER\n")
	_, err := ParseFile(path)
	if err == nil {
		t.Fatal("expected error for malformed line, got nil")
	}
}

func TestParseFile_FileNotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
