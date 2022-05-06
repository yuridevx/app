package apptrace

import "context"

type traceKey int

var traceKeyVal traceKey

func TraceContext(ctx context.Context, trace *Trace) context.Context {
	return context.WithValue(ctx, traceKeyVal, trace)
}

func TraceContextNew(ctx context.Context) (context.Context, *Trace) {
	trace := NewTrace()
	return TraceContext(ctx, trace), trace
}

func FromContext(ctx context.Context) *Trace {
	trace, _ := ctx.Value(traceKeyVal).(*Trace)
	return trace
}
