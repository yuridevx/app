package options

import (
	"time"
)

type CPeriodOptions struct {
	Handler       HandlerFn
	SwitchCh      chan bool
	Period        time.Duration
	PeriodFn      PeriodFn
	PeriodIndexFn PeriodIndexFn
	Middleware    []Middleware
}

type CPeriodOption func(o *CPeriodOptions)

func DefaultCPeriodOptions() CPeriodOptions {
	return CPeriodOptions{}
}

func (a *CPeriodOptions) Validate() {
	if a.Handler == nil {
		panic("Handler is required")
	}
	if a.PeriodFn == nil && a.PeriodIndexFn == nil && a.Period == 0 {
		panic("Period is required")
	}
}
