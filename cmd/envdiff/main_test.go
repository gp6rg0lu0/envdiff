package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return path
}

func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "envdiff")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func TestMain_NoDifferences(t *testing.T) {
	bin := buildBinary(t)
	base := writeTempEnv(t, "KEY=value\nFOO=bar\n")
	cmp := writeTempEnv(t, "KEY=value\nFOO=bar\n")

	cmd := exec.Command(bin, base, cmp)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got error: %v\noutput: %s", err, out)
	}
}

func TestMain_MissingKey_ExitsTwo(t *testing.T) {
	bin := buildBinary(t)
	base := writeTempEnv(t, "KEY=value\nMISSING=x\n")
	cmp := writeTempEnv(t, "KEY=value\n")

	cmd := exec.Command(bin, base, cmp)
	out, err := cmd.CombinedOutput()
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 2 {
			t.Fatalf("expected exit code 2, got %d\noutput: %s", exitErr.ExitCode(), out)
		}
	} else if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else {
		t.Fatalf("expected exit code 2 but got 0\noutput: %s", out)
	}
}

func TestMain_JSONFormat(t *testing.T) {
	bin := buildBinary(t)
	base := writeTempEnv(t, "KEY=value\n")
	cmp := writeTempEnv(t, "KEY=value\n")

	cmd := exec.Command(bin, "-format", "json", base, cmp)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got error: %v\noutput: %s", err, out)
	}
	if len(out) == 0 {
		t.Fatal("expected JSON output, got empty")
	}
}

func TestMain_MissingArgs(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	if err := cmd.Run(); err == nil {
		t.Fatal("expected non-zero exit for missing args")
	}
}
