package apprelic

import (
	"context"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/yuridevx/app/apptrace"
	"github.com/yuridevx/app/options"
)

func RelicTrace(txn *newrelic.Transaction, trace *apptrace.Trace, err error) {
	if txn == nil || trace == nil {
		return
	}

	for k, v := range trace.GetAttributes() {
		txn.AddAttribute(k, v)
	}

	if err != nil {
		txn.NoticeError(err)
	}

	for _, noticeErr := range trace.GetErrors() {
		txn.NoticeError(noticeErr)
	}
}

func NewRelicTransactionMiddleware(
	nr *newrelic.Application,
) options.Middleware {
	return func(ctx context.Context, input interface{}, call options.Call, next options.NextFn) error {
		// prepare context
		txn := nr.StartTransaction("")
		defer txn.End()
		ctx = newrelic.NewContext(ctx, txn)
		// execute
		err := next(ctx, input)
		return err
	}
}

func NewRelicTraceMiddleware() options.Middleware {
	return func(
		ctx context.Context,
		input interface{},
		part options.Call,
		next options.NextFn,
	) error {
		// prepare context
		trace := apptrace.NewTrace()
		ctx = apptrace.TraceContext(ctx, trace)

		// execute
		err := next(ctx, input)
		txn := newrelic.FromContext(ctx)
		if txn != nil {
			txn.SetName(trace.GetName())
			RelicTrace(txn, trace, err)
		}

		return err
	}
}
