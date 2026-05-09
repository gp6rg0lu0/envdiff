package diff

import (
	"fmt"
	"testing"
)

func BenchmarkAnnotateResult_Large(b *testing.B) {
	r := Result{}
	for i := 0; i < 100; i++ {
		r.Missing = append(r.Missing, Entry{
			Key:       fmt.Sprintf("MISSING_KEY_%d", i),
			BaseValue: fmt.Sprintf("value%d", i),
		})
		r.Extra = append(r.Extra, Entry{
			Key:         fmt.Sprintf("EXTRA_KEY_%d", i),
			TargetValue: fmt.Sprintf("value%d", i),
		})
		r.Mismatched = append(r.Mismatched, Entry{
			Key:         fmt.Sprintf("MISMATCH_KEY_%d", i),
			BaseValue:   fmt.Sprintf("base%d", i),
			TargetValue: fmt.Sprintf("target%d", i),
		})
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		AnnotateResult(r)
	}
}

func BenchmarkAnnotateResult_SensitiveKeys(b *testing.B) {
	r := Result{}
	for i := 0; i < 50; i++ {
		r.Extra = append(r.Extra, Entry{
			Key:         fmt.Sprintf("SECRET_KEY_%d", i),
			TargetValue: fmt.Sprintf("secret%d", i),
		})
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		AnnotateResult(r)
	}
}
