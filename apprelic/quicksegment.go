package apptrace

import (
	"context"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/yuridevx/app/apptrace"
	"runtime"
)

func noop() {
}

type SegmentOption func(segment *newrelic.Segment)

func WithName(name string) SegmentOption {
	return func(segment *newrelic.Segment) {
		segment.Name = name
	}
}

func WithAttributesMap(attributes map[string]interface{}) SegmentOption {
	return func(segment *newrelic.Segment) {
		for k, v := range attributes {
			segment.AddAttribute(k, v)
		}
	}
}

func WithAttributes(attributes ...interface{}) SegmentOption {
	return func(segment *newrelic.Segment) {
		for i := 0; i < len(attributes); i += 2 {
			segment.AddAttribute(attributes[i].(string), attributes[i+1])
		}
	}
}

/*
QuickSegment creates a new segment with the given name and options.
```go
func SomewhereInTheCode(id string) {
	defer QuickSegment(apptrace.WithAttributes("id", id)()
	// ...
}
```
*/
func QuickSegment(ctx context.Context, options ...SegmentOption) func() {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return noop
	}
	segment := txn.StartSegment("")
	for _, option := range options {
		option(segment)
	}
	if segment.Name == "" {
		caller, _, _, _ := runtime.Caller(1)
		segment.Name = apptrace.FunctionToName(caller)
	}
	return segment.End
}
