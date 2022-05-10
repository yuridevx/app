package apptrace

import (
	"context"
	"github.com/yuridevx/app/options"
)

func NewTraceMiddleware(
	naming NamingStrategy,
) options.Middleware {
	if naming == nil {
		naming = DefaultNamingStrategy
	}
	return func(
		ctx context.Context,
		input interface{},
		part options.Call,
		next options.NextFn,
	) error {
		if part == nil {
			err := next(ctx, input)
			return err
		}
		name := naming(part)
		trace := NewTrace().WithName(name)
		ctx = TraceContext(ctx, trace)
		err := next(ctx, input)
		return err
	}
}
