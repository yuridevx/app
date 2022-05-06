package apptrace

import (
	"context"
	"github.com/yuridevx/app/extension"
)

func NewTraceMiddleware(
	naming NamingStrategy,
) extension.Middleware {
	if naming == nil {
		naming = DefaultNamingStrategy
	}
	return func(
		ctx context.Context,
		call extension.CallType,
		input interface{},
		part extension.Part,
		next extension.NextFn,
	) error {
		name := naming(call, part)
		trace := NewTrace().WithName(name)
		ctx = TraceContext(ctx, trace)
		err := next(ctx, input)
		return err
	}
}
