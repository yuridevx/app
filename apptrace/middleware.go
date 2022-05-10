package apptrace

import (
	"context"
	"github.com/yuridevx/app/options"
)

func NewTraceMiddleware(
	naming NamingStrategy,
	ignoreCallTypes []options.CallType,
) options.Middleware {
	if naming == nil {
		naming = DefaultNamingStrategy
	}
	if ignoreCallTypes == nil {
		ignoreCallTypes = []options.CallType{
			options.CallStartGroup,
			options.CallShutdownGroup,
		}
	}
	return func(
		ctx context.Context,
		input interface{},
		part options.Call,
		next options.NextFn,
	) error {
		if part.IsCallType(ignoreCallTypes...) {
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
