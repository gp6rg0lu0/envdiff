package diff

// AnnotationLevel indicates the severity of an annotation.
type AnnotationLevel string

const (
	AnnotationInfo    AnnotationLevel = "info"
	AnnotationWarning AnnotationLevel = "warning"
	AnnotationError   AnnotationLevel = "error"
)

// Annotation attaches a human-readable note to a specific key.
type Annotation struct {
	Key     string
	Level   AnnotationLevel
	Message string
}

// AnnotateResult inspects a Result and produces a slice of Annotations
// describing notable findings such as empty values, sensitive keys exposed
// as extra, or keys present in one environment but not the other.
func AnnotateResult(r Result) []Annotation {
	var annotations []Annotation

	for _, e := range r.Missing {
		annotations = append(annotations, Annotation{
			Key:     e.Key,
			Level:   AnnotationError,
			Message: "key is missing from target environment",
		})
	}

	for _, e := range r.Extra {
		level := AnnotationWarning
		msg := "key exists only in target environment"
		if IsSensitive(e.Key) {
			level = AnnotationError
			msg = "sensitive key exists only in target environment — possible secret leak"
		}
		annotations = append(annotations, Annotation{
			Key:     e.Key,
			Level:   level,
			Message: msg,
		})
	}

	for _, e := range r.Mismatched {
		level := AnnotationWarning
		msg := "value differs between environments"
		if e.BaseValue == "" || e.TargetValue == "" {
			level = AnnotationError
			msg = "value differs and one side is empty"
		}
		annotations = append(annotations, Annotation{
			Key:     e.Key,
			Level:   level,
			Message: msg,
		})
	}

	return annotations
}
