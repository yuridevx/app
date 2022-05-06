package options

import (
	"github.com/yuridevx/app/extension"
	"reflect"
)

type CConsumeOptions struct {
	Handler        HandlerFn
	ConsumeCH      interface{}
	CallMiddleware []extension.Middleware
}

func DefaultCConsumeOptions() CConsumeOptions {
	return CConsumeOptions{}
}

func (a *CConsumeOptions) Merge(from *CConsumeOptions) {
	if from.Handler != nil {
		a.Handler = from.Handler
	}
	if from.ConsumeCH != nil {
		a.ConsumeCH = from.ConsumeCH
	}
	if from.CallMiddleware != nil {
		merged := make([]extension.Middleware, len(a.CallMiddleware)+len(from.CallMiddleware))
		copy(merged, a.CallMiddleware)
		copy(merged[len(a.CallMiddleware):], from.CallMiddleware)
		a.CallMiddleware = merged
	}
}

func (a *CConsumeOptions) Validate() {
	if a.Handler == nil {
		panic("Handler is required")
	}
	if reflect.TypeOf(a.Handler).Kind() != reflect.Func {
		panic("Handler must be a function")
	}
	if a.ConsumeCH == nil {
		panic("ConsumeCH is required")
	}
	if reflect.TypeOf(a.ConsumeCH).Kind() != reflect.Chan {
		panic("ConsumeCH must be a channel")
	}
}
