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

type ShouldLog = func(
	trace *apptrace.Trace,
	callType extension.CallType,
	time LogTime,
) bool

var ShouldLogAfter = func(trace *apptrace.Trace, callType extension.CallType, time LogTime) bool {
	return trace.Log && time == LogAfter
}

var ShouldLogBefore = func(trace *apptrace.Trace, callType extension.CallType, time LogTime) bool {
	return trace.Log && time == LogBefore
}

var ShouldLogAll = func(trace *apptrace.Trace, callType extension.CallType, time LogTime) bool {
	return trace.Log
}

var ShouldLogAlways = func(trace *apptrace.Trace, callType extension.CallType, time LogTime) bool {
	return true
}

var ShouldLogNever = func(trace *apptrace.Trace, callType extension.CallType, time LogTime) bool {
	return false
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
	shouldLog ShouldLog,
) extension.Middleware {
	if shouldLog == nil {
		shouldLog = ShouldLogAll
	}
	return func(
		ctx context.Context,
		call extension.CallType,
		input interface{},
		part extension.Part,
		next extension.NextFn,
	) error {
		trace := apptrace.FromContext(ctx)
		if shouldLog(trace, call, LogBefore) {
			logger.Info(trace.GetName() + " starting")
		}
		err := next(ctx, input)
		if shouldLog(trace, call, LogAfter) {
			if err != nil {
				logger.Error(trace.GetName()+" finished", AppErr(err), InlineMap(trace.GetAttributes()))
			} else {
				logger.Info(trace.GetName()+" finished", InlineMap(trace.GetAttributes()))
			}
			for _, noticeErr := range trace.GetErrors() {
				logger.Warn(trace.GetName()+" notice ", AppErr(noticeErr))
			}
		}
		return err
	}
}
