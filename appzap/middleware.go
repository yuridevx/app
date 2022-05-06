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
		prefix := string(call) + " " + apptrace.FunctionToName(part.GetHandler())
		if shouldLog(trace, call, LogBefore) {
			logger.Info(prefix + " starting")
		}
		err := next(ctx, input)
		if shouldLog(trace, call, LogAfter) {
			if err != nil {
				logger.Error(prefix+" finished", AppErr(err), InlineMap(trace.GetAttributes()))
			} else {
				logger.Info(prefix+" finished", InlineMap(trace.GetAttributes()))
			}
			for _, noticeErr := range trace.GetErrors() {
				logger.Warn(prefix+" notice ", AppErr(noticeErr))
			}
		}
		return err
	}
}
