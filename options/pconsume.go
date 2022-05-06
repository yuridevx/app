package options

import (
	"github.com/yuridevx/app/extension"
	"reflect"
)

type PConsumeOptions struct {
	Handler        HandlerFn
	ConsumeCH      interface{}
	Goroutines     int
	CallMiddleware []extension.Middleware
}

func DefaultPConsumeOptions() PConsumeOptions {
	return PConsumeOptions{
		Goroutines: 1,
	}
}

func (a *PConsumeOptions) Merge(from *PConsumeOptions) {
	if from.Handler != nil {
		a.Handler = from.Handler
	}
	if from.ConsumeCH != nil {
		a.ConsumeCH = from.ConsumeCH
	}
	if from.Goroutines != 0 {
		a.Goroutines = from.Goroutines
	}
	if from.CallMiddleware != nil {
		merged := make([]extension.Middleware, len(a.CallMiddleware)+len(from.CallMiddleware))
		copy(merged, a.CallMiddleware)
		copy(merged[len(a.CallMiddleware):], from.CallMiddleware)
		a.CallMiddleware = merged
	}
}

func (a *PConsumeOptions) Validate() {
	if a.Handler == nil {
		panic("handler is required")
	}
	if reflect.TypeOf(a.Handler).Kind() != reflect.Func {
		panic("handler must be a function")
	}
	if a.ConsumeCH == nil {
		panic("consume ch is required")
	}
	if reflect.TypeOf(a.ConsumeCH).Kind() != reflect.Chan {
		panic("consume ch must be a channel")
	}
	if a.Goroutines <= 0 {
		panic("goroutines must be greater than 0")
	}
}
