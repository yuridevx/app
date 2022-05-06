package apptrace

import (
	"context"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/yuridevx/app/extension"
	"go.uber.org/zap"
	"sync"
)

var DefaultLogTypes = []extension.CallType{
	extension.CallPeriodic,
	extension.CallCConsume,
	extension.CallStart,
	extension.CallShutdown,
}

//NewNewRelicMiddleware use it in app to add newrelic tracing and zap logging without code changes
func NewNewRelicMiddleware(
	logger *zap.Logger,
	application *newrelic.Application,
	logCallTypes []extension.CallType,
) extension.Middleware {
	if logCallTypes == nil {
		logCallTypes = DefaultLogTypes
	}
	return func(
		ctx context.Context,
		call extension.CallType,
		input interface{},
		part extension.Part,
		next extension.NextFn,
	) error {
		// prepare context
		handler := part.GetHandler()
		name := FunctionToName(handler)
		txn := application.StartTransaction(name)
		defer txn.End()
		ctx = newrelic.NewContext(ctx, txn)
		sm := &sync.Map{}
		ctx = AttributesContext(ctx, sm)

		// execute
		err := next(ctx, input)

		// add attributes to transaction
		hasData := false
		sm.Range(func(key, value interface{}) bool {
			keyS, ok := key.(string)
			if !ok {
				return true
			}
			hasData = true
			txn.AddAttribute(keyS, value)
			return true
		})

		// if attributes are not empty log it
		if hasData && contains(logCallTypes, call) {
			if err == nil {
				logger.Info(name, InlineSyncMap(sm))
			} else {
				logger.Error(name+" error", zap.Error(err), InlineSyncMap(sm))
			}
		}
		if err != nil {
			txn.NoticeError(err)
		}
		return err
	}
}

func contains[E comparable](s []E, v E) bool {
	for _, vs := range s {
		if v == vs {
			return true
		}
	}
	return false
}
