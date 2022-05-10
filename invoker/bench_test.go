package invoker

import (
	"context"
	"testing"
)

type someType struct {
}

func BenchmarkCast(b *testing.B) {
	plain := NewInvoker(func() {})

	b.Run("Plain", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = plain.Invoke(context.Background(), nil, nil, nil)
		}
	})

	cast := NewInvoker(func(input *someType) {})

	b.Run("Cast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = cast.Invoke(context.Background(), nil, nil, nil)
		}
	})
}
