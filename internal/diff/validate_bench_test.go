package diff

import (
	"fmt"
	"testing"
)

func BenchmarkValidate_LargeEnv(b *testing.B) {
	env := make(map[string]string, 500)
	required := make([]string, 50)
	forbidden := make([]string, 10)

	for i := 0; i < 500; i++ {
		env[fmt.Sprintf("KEY_%d", i)] = fmt.Sprintf("value_%d", i)
	}
	for i := 0; i < 50; i++ {
		required[i] = fmt.Sprintf("KEY_%d", i)
	}
	for i := 0; i < 10; i++ {
		forbidden[i] = fmt.Sprintf("FORBIDDEN_%d", i)
	}

	opts := ValidateOptions{
		DisallowEmptyValues: true,
		RequiredKeys:        required,
		ForbiddenKeys:       forbidden,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Validate(env, opts)
	}
}

func BenchmarkValidate_EmptyValues(b *testing.B) {
	env := make(map[string]string, 200)
	for i := 0; i < 100; i++ {
		env[fmt.Sprintf("FILLED_%d", i)] = fmt.Sprintf("val_%d", i)
	}
	for i := 0; i < 100; i++ {
		env[fmt.Sprintf("EMPTY_%d", i)] = ""
	}
	opts := ValidateOptions{DisallowEmptyValues: true}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Validate(env, opts)
	}
}
