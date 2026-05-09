package diff

import (
	"testing"
)

func TestAnnotateResult_NoIssues(t *testing.T) {
	r := Result{}
	anns := AnnotateResult(r)
	if len(anns) != 0 {
		t.Fatalf("expected 0 annotations, got %d", len(anns))
	}
}

func TestAnnotateResult_MissingKey(t *testing.T) {
	r := Result{
		Missing: []Entry{{Key: "DB_HOST", BaseValue: "localhost"}},
	}
	anns := AnnotateResult(r)
	if len(anns) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(anns))
	}
	if anns[0].Level != AnnotationError {
		t.Errorf("expected error level, got %s", anns[0].Level)
	}
	if anns[0].Key != "DB_HOST" {
		t.Errorf("expected key DB_HOST, got %s", anns[0].Key)
	}
}

func TestAnnotateResult_ExtraKey_NonSensitive(t *testing.T) {
	r := Result{
		Extra: []Entry{{Key: "APP_NAME", TargetValue: "myapp"}},
	}
	anns := AnnotateResult(r)
	if len(anns) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(anns))
	}
	if anns[0].Level != AnnotationWarning {
		t.Errorf("expected warning level for non-sensitive extra key, got %s", anns[0].Level)
	}
}

func TestAnnotateResult_ExtraKey_Sensitive(t *testing.T) {
	r := Result{
		Extra: []Entry{{Key: "SECRET_TOKEN", TargetValue: "abc123"}},
	}
	anns := AnnotateResult(r)
	if len(anns) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(anns))
	}
	if anns[0].Level != AnnotationError {
		t.Errorf("expected error level for sensitive extra key, got %s", anns[0].Level)
	}
}

func TestAnnotateResult_MismatchedValue_EmptySide(t *testing.T) {
	r := Result{
		Mismatched: []Entry{{Key: "PORT", BaseValue: "8080", TargetValue: ""}},
	}
	anns := AnnotateResult(r)
	if len(anns) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(anns))
	}
	if anns[0].Level != AnnotationError {
		t.Errorf("expected error level when one side is empty, got %s", anns[0].Level)
	}
}

func TestAnnotateResult_MismatchedValue_BothPresent(t *testing.T) {
	r := Result{
		Mismatched: []Entry{{Key: "LOG_LEVEL", BaseValue: "debug", TargetValue: "info"}},
	}
	anns := AnnotateResult(r)
	if len(anns) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(anns))
	}
	if anns[0].Level != AnnotationWarning {
		t.Errorf("expected warning level for value mismatch, got %s", anns[0].Level)
	}
}

func TestAnnotateResult_MultipleAnnotations(t *testing.T) {
	r := Result{
		Missing:    []Entry{{Key: "A", BaseValue: "1"}},
		Extra:      []Entry{{Key: "B", TargetValue: "2"}},
		Mismatched: []Entry{{Key: "C", BaseValue: "x", TargetValue: "y"}},
	}
	anns := AnnotateResult(r)
	if len(anns) != 3 {
		t.Fatalf("expected 3 annotations, got %d", len(anns))
	}
}
