package options

import (
	"time"
)

type ApplicationOptions struct {
	StartTimeout             time.Duration
	AppShutdownTimeout       time.Duration
	ComponentShutdownTimeout time.Duration
	Middleware               []Middleware
}

type ApplicationOption func(o *ApplicationOptions)

func DefaultApplicationOptions() ApplicationOptions {
	return ApplicationOptions{
		StartTimeout:             time.Second * 60,
		AppShutdownTimeout:       time.Second * 20,
		ComponentShutdownTimeout: time.Second * 40,
	}
}

func (a ApplicationOptions) Validate() {
	if a.StartTimeout < time.Second {
		panic("StartTimeout must be greater than or equal to 1 second")
	}
	if a.AppShutdownTimeout < time.Second {
		panic("AppShutdownTimeout must be greater than or equal to 1 second")
	}
	if a.ComponentShutdownTimeout < time.Second {
		panic("ComponentShutdownTimeout must be greater than or equal to 1 second")
	}
}
