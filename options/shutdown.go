package options

import (
	"github.com/yuridevx/app/extension"
	"reflect"
)

type ComponentShutdownOptions struct {
	Handler        HandlerFn
	CallMiddleware []extension.Middleware
}

func DefaultComponentShutdownOptions() ComponentShutdownOptions {
	return ComponentShutdownOptions{}
}

func (a *ComponentShutdownOptions) Merge(from *ComponentShutdownOptions) {
	if from.Handler != nil {
		a.Handler = from.Handler
	}
	if from.CallMiddleware != nil {
		merged := make([]extension.Middleware, len(a.CallMiddleware)+len(from.CallMiddleware))
		copy(merged, a.CallMiddleware)
		copy(merged[len(a.CallMiddleware):], from.CallMiddleware)
		a.CallMiddleware = merged
	}
}

func (a *ComponentShutdownOptions) Validate() {
	if a.Handler == nil {
		panic("handler is required")
	}
	if reflect.TypeOf(a.Handler).Kind() != reflect.Func {
		panic("handler must be a function")
	}
}
