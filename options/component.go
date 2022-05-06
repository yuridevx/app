package options

import (
	"github.com/yuridevx/app/extension"
)

type ComponentOptions struct {
	Definition          ComponentDefinition
	ComponentMiddleware []extension.Middleware
}

func DefaultComponentOptions() ComponentOptions {
	return ComponentOptions{}
}

func (a *ComponentOptions) Merge(from *ComponentOptions) {
	if from.Definition != nil {
		a.Definition = from.Definition
	}
	if from.ComponentMiddleware != nil {
		merged := make([]extension.Middleware, len(a.ComponentMiddleware)+len(from.ComponentMiddleware))
		copy(merged, a.ComponentMiddleware)
		copy(merged[len(a.ComponentMiddleware):], from.ComponentMiddleware)
		a.ComponentMiddleware = merged
	}
}
