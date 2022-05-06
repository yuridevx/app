package appzap

import (
	"context"
	"github.com/yuridevx/app/apptrace"
	"github.com/yuridevx/app/extension"
	"go.uber.org/zap"
)

type LogTime string

const LogBefore LogTime = "before"
const LogAfter LogTime = "after"

type LogFn = func(
	trace *apptrace.Trace,
	callType extension.CallType,
	time LogTime,
) (bool, []zap.Field)

var DefaultLogFn = func(
	trace *apptrace.Trace,
	callType extension.CallType,
	time LogTime,
) (bool, []zap.Field) {
	if callType == extension.CallPBlocking ||
		callType == extension.CallStart ||
		callType == extension.CallShutdown {
		return true, nil
	}
	return trace.GetLog(), nil
}

var LogMeMiddleware = func(
	ctx context.Context,
	call extension.CallType,
	input interface{},
	part extension.Part,
	next extension.NextFn,
) error {
	apptrace.FromContext(ctx).WithLog(true)
	return next(ctx, input)
}

var DontLogMeMiddleware = func(
	ctx context.Context,
	call extension.CallType,
	input interface{},
	part extension.Part,
	next extension.NextFn,
) error {
	apptrace.FromContext(ctx).WithLog(false)
	return next(ctx, input)
}

func ZapMiddleware(
	logger *zap.Logger,
	shouldLog LogFn,
) extension.Middleware {
	if shouldLog == nil {
		shouldLog = DefaultLogFn
	}
	return func(
		ctx context.Context,
		call extension.CallType,
		input interface{},
		part extension.Part,
		next extension.NextFn,
	) error {
		trace := apptrace.FromContext(ctx)
		log, fields := shouldLog(trace, call, LogBefore)
		if log {
			logger.Info(trace.GetName()+" starting", fields...)
		}
		err := next(ctx, input)
		log, fields = shouldLog(trace, call, LogAfter)
		if log {
			if err != nil {
				logger.With(AppErr(err), InlineMap(trace.GetAttributes())).Error(trace.GetName()+" finished", fields...)
			} else {
				logger.With(InlineMap(trace.GetAttributes())).Info(trace.GetName()+" finished", fields...)
			}
			for _, noticeErr := range trace.GetErrors() {
				logger.Warn(trace.GetName()+" notice ", AppErr(noticeErr))
			}
		}
		return err
	}
}
