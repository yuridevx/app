package options

import (
	"github.com/yuridevx/app/extension"
	"reflect"
	"time"
)

type CPeriodOptions struct {
	Handler        HandlerFn
	SwitchCh       chan bool
	Period         time.Duration
	PeriodFn       PeriodFn
	PeriodIndexFn  PeriodIndexFn
	CallMiddleware []extension.Middleware
}

func DefaultCPeriodOptions() CPeriodOptions {
	return CPeriodOptions{}
}

func (a *CPeriodOptions) Merge(from *CPeriodOptions) {
	if from.Handler != nil {
		a.Handler = from.Handler
	}
	if from.SwitchCh != nil {
		a.SwitchCh = from.SwitchCh
	}
	if from.Period != 0 {
		a.Period = from.Period
	}
	if from.PeriodFn != nil {
		a.PeriodFn = from.PeriodFn
	}
	if from.PeriodIndexFn != nil {
		a.PeriodIndexFn = from.PeriodIndexFn
	}
	if from.CallMiddleware != nil {
		merged := make([]extension.Middleware, len(a.CallMiddleware)+len(from.CallMiddleware))
		copy(merged, a.CallMiddleware)
		copy(merged[len(a.CallMiddleware):], from.CallMiddleware)
		a.CallMiddleware = merged
	}
}

func (a *CPeriodOptions) Validate() {
	if a.Handler == nil {
		panic("Handler is required")
	}
	if reflect.TypeOf(a.Handler).Kind() != reflect.Func {
		panic("Handler must be a function")
	}
	if a.Period != 0 {
		return
	}
	if a.PeriodFn != nil {
		if reflect.TypeOf(a.PeriodFn).Kind() != reflect.Func {
			panic("PeriodFn must be a function")
		}
		return
	}
	if a.PeriodIndexFn != nil {
		if reflect.TypeOf(a.PeriodIndexFn).Kind() != reflect.Func {
			panic("PeriodIndexFn must be a function")
		}
		return
	}
}
