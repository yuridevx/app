package options

import (
	"github.com/yuridevx/app/extension"
	"time"
)

type ApplicationOptions struct {
	StartTimeout             time.Duration
	AppShutdownTimeout       time.Duration
	ComponentShutdownTimeout time.Duration
	GlobalMiddleware         []extension.Middleware
}

func DefaultApplicationOptions() ApplicationOptions {
	return ApplicationOptions{
		StartTimeout:             time.Second * 60,
		AppShutdownTimeout:       time.Second * 20,
		ComponentShutdownTimeout: time.Second * 40,
	}
}

func (a *ApplicationOptions) Merge(from *ApplicationOptions) {
	if from.StartTimeout != 0 {
		a.StartTimeout = from.StartTimeout
	}
	if from.AppShutdownTimeout != 0 {
		a.AppShutdownTimeout = from.AppShutdownTimeout
	}
	if from.GlobalMiddleware != nil {
		merged := make([]extension.Middleware, len(a.GlobalMiddleware)+len(from.GlobalMiddleware))
		copy(merged, a.GlobalMiddleware)
		copy(merged[len(a.GlobalMiddleware):], from.GlobalMiddleware)
		a.GlobalMiddleware = merged
	}
}
