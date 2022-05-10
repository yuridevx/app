package apptrace

import "context"

type traceKey int

var traceKeyVal traceKey

const GinAppTraceContextKey = "GinAppTraceContextKey" // used for gin context

func TraceContext(ctx context.Context, trace *Trace) context.Context {
	// gin support
	if ginCtx, ok := ctx.(interface {
		Set(key string, value interface{})
	}); ok {
		ginCtx.Set(GinAppTraceContextKey, trace)
		return ctx
	}

	return context.WithValue(ctx, traceKeyVal, trace)
}

func TraceContextNew(ctx context.Context) (context.Context, *Trace) {
	trace := NewTrace()
	return TraceContext(ctx, trace), trace
}

func FromContext(ctx context.Context) *Trace {
	trace, _ := ctx.Value(traceKeyVal).(*Trace)
	if trace != nil {
		return trace
	}
	trace, _ = ctx.Value(GinAppTraceContextKey).(*Trace)
	if trace != nil {
		return trace
	}
	return nil
}
