package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// RenderAnnotations renders a list of annotations in the requested format.
func RenderAnnotations(annotations []diff.Annotation, format string) string {
	switch strings.ToLower(format) {
	case "json":
		return renderAnnotationsJSON(annotations)
	default:
		return renderAnnotationsText(annotations)
	}
}

func renderAnnotationsText(annotations []diff.Annotation) string {
	if len(annotations) == 0 {
		return "No annotations.\n"
	}

	var sb strings.Builder
	sb.WriteString("Annotations:\n")
	for _, a := range annotations {
		icon := levelIcon(a.Level)
		fmt.Fprintf(&sb, "  %s [%s] %s: %s\n", icon, strings.ToUpper(string(a.Level)), a.Key, a.Message)
	}
	return sb.String()
}

func levelIcon(level diff.AnnotationLevel) string {
	switch level {
	case diff.AnnotationError:
		return "✖"
	case diff.AnnotationWarning:
		return "⚠"
	default:
		return "ℹ"
	}
}

type annotationJSON struct {
	Key     string `json:"key"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

func renderAnnotationsJSON(annotations []diff.Annotation) string {
	type payload struct {
		Annotations []annotationJSON `json:"annotations"`
		Count       int              `json:"count"`
	}

	items := make([]annotationJSON, len(annotations))
	for i, a := range annotations {
		items[i] = annotationJSON{
			Key:     a.Key,
			Level:   string(a.Level),
			Message: a.Message,
		}
	}

	out, _ := json.MarshalIndent(payload{Annotations: items, Count: len(items)}, "", "  ")
	return string(out) + "\n"
}
