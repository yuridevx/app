package apptrace

import (
	"context"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/yuridevx/app/apptrace"
	"github.com/yuridevx/app/extension"
)

func NewNewRelicMiddleware(
	application *newrelic.Application,
) extension.Middleware {
	return func(
		ctx context.Context,
		call extension.CallType,
		input interface{},
		part extension.Part,
		next extension.NextFn,
	) error {
		// prepare context
		trace := apptrace.FromContext(ctx)
		handler := part.GetHandler()
		name := apptrace.FunctionToName(handler)
		txn := application.StartTransaction(name)
		defer txn.End()
		ctx = newrelic.NewContext(ctx, txn)

		// execute
		err := next(ctx, input)

		for k, v := range trace.GetAttributes() {
			txn.AddAttribute(k, v)
		}

		if err != nil {
			txn.NoticeError(err)
		}

		for _, noticeErr := range trace.GetErrors() {
			txn.NoticeError(noticeErr)
		}

		return err
	}
}
