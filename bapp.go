package app

import (
	"github.com/yuridevx/app/handlers"
	"github.com/yuridevx/app/options"
)

type appBuilder struct {
	appOpts    options.ApplicationOptions
	events     handlers.Events
	components []*componentBuilder
	shutdown   []ShutdownFn
}

func (a *appBuilder) Events(events handlers.Events) Builder {
	a.events = events
	return a
}

func (a *appBuilder) Options(opts ...options.ApplicationOption) Builder {
	for _, o := range opts {
		o(&a.appOpts)
	}
	return a
}

func (a *appBuilder) C(def options.ComponentDefinition, opts ...options.ComponentOption) ComponentBuilder {
	c := newComponentBuilder(def, opts...)
	a.components = append(a.components, c)
	return c
}

func (a *appBuilder) OnShutdown(fn ...ShutdownFn) Builder {
	a.shutdown = append(a.shutdown, fn...)
	return a
}

func (a *appBuilder) Build() Application {
	ap := &app{
		opts:           a.appOpts,
		components:     make([]*component, len(a.components)),
		shutdown:       make([]ShutdownFn, len(a.shutdown)),
		shutdownSignal: make(chan struct{}),
	}
	for i, c := range a.components {
		ap.components[i] = c.build(a.appOpts, a.events)
	}
	for i, s := range a.shutdown {
		ap.shutdown[i] = s
	}
	a.reset()
	return ap
}

func (a *appBuilder) reset() {
	a.appOpts = options.ApplicationOptions{}
	a.components = nil
	a.shutdown = nil
}

func NewBuilder(opts ...options.ApplicationOption) Builder {
	opt := options.DefaultApplicationOptions()
	for _, o := range opts {
		o(&opt)
	}
	return &appBuilder{
		events:  handlers.NullEvents{},
		appOpts: opt,
	}
}
