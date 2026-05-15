package diff

import (
	"strings"
	"testing"
)

func TestMaskValue_EmptyString(t *testing.T) {
	opts := DefaultMaskOptions()
	if got := MaskValue("", opts); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestMaskValue_ShortValue_FullyMasked(t *testing.T) {
	opts := DefaultMaskOptions() // MinLength=6
	got := MaskValue("abc", opts)
	if got != "***" {
		t.Errorf("expected \"***\", got %q", got)
	}
}

func TestMaskValue_LongValue_PartiallyMasked(t *testing.T) {
	opts := DefaultMaskOptions() // VisibleChars=3
	got := MaskValue("supersecret", opts)
	if !strings.HasPrefix(got, "sup") {
		t.Errorf("expected prefix \"sup\", got %q", got)
	}
	if !strings.Contains(got[3:], "*") {
		t.Errorf("expected mask chars after prefix, got %q", got)
	}
	if len(got) != len("supersecret") {
		t.Errorf("expected same length as input, got %d", len(got))
	}
}

func TestMaskValue_CustomMaskChar(t *testing.T) {
	opts := MaskOptions{VisibleChars: 2, MaskChar: '#', MinLength: 4}
	got := MaskValue("password", opts)
	if !strings.HasPrefix(got, "pa") {
		t.Errorf("expected prefix \"pa\", got %q", got)
	}
	if !strings.Contains(got, "#") {
		t.Errorf("expected '#' mask char, got %q", got)
	}
}

func TestMaskValue_ExactMinLength(t *testing.T) {
	opts := DefaultMaskOptions() // MinLength=6
	got := MaskValue("abcdef", opts)
	// length == MinLength, should partially reveal
	if !strings.HasPrefix(got, "abc") {
		t.Errorf("expected prefix \"abc\", got %q", got)
	}
}

func TestMaskResult_MasksSensitiveValues(t *testing.T) {
	result := Result{
		Entries: []Entry{
			{Key: "API_SECRET", Status: StatusMismatched, BaseValue: "oldsecret123", OtherValue: "newsecret456"},
			{Key: "APP_NAME", Status: StatusMismatched, BaseValue: "myapp", OtherValue: "otherapp"},
		},
	}
	opts := DefaultMaskOptions()
	masked := MaskResult(result, opts)

	secretEntry := masked.Entries[0]
	if secretEntry.BaseValue == "oldsecret123" {
		t.Error("expected BaseValue to be masked for API_SECRET")
	}
	if secretEntry.OtherValue == "newsecret456" {
		t.Error("expected OtherValue to be masked for API_SECRET")
	}

	appEntry := masked.Entries[1]
	if appEntry.BaseValue != "myapp" {
		t.Errorf("expected BaseValue unchanged for APP_NAME, got %q", appEntry.BaseValue)
	}
}

func TestMaskResult_DoesNotMutateOriginal(t *testing.T) {
	original := Result{
		Entries: []Entry{
			{Key: "DB_PASSWORD", Status: StatusExtra, OtherValue: "hunter2hunter"},
		},
	}
	opts := DefaultMaskOptions()
	_ = MaskResult(original, opts)

	if original.Entries[0].OtherValue != "hunter2hunter" {
		t.Error("MaskResult mutated the original result")
	}
}
