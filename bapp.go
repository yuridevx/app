package app

import (
	"github.com/yuridevx/app/handlers"
	"github.com/yuridevx/app/options"
)

type appBuilder struct {
	options.ApplicationOptions
	events     handlers.Events
	components []*componentBuilder
	shutdown   []AppShutdownFn
}

func (a *appBuilder) Events(events handlers.Events) AppBuilder {
	a.events = events
	return a
}

func (a *appBuilder) Options(opts options.ApplicationOptions) AppBuilder {
	a.ApplicationOptions.Merge(&opts)
	return a
}

func (a *appBuilder) C(def options.ComponentDefinition, opts ...options.ComponentOptions) ComponentBuilder {
	c := newComponentBuilder(def, opts...)
	a.components = append(a.components, c)
	return c
}

func (a *appBuilder) OnShutdown(fn ...AppShutdownFn) AppBuilder {
	a.shutdown = append(a.shutdown, fn...)
	return a
}

func (a *appBuilder) Build() Application {
	ap := &app{
		opts:           a.ApplicationOptions,
		components:     make([]*component, len(a.components)),
		shutdown:       make([]AppShutdownFn, len(a.shutdown)),
		shutdownSignal: make(chan struct{}),
	}
	for i, c := range a.components {
		ap.components[i] = c.build(a.ApplicationOptions, a.events)
	}
	for i, s := range a.shutdown {
		ap.shutdown[i] = s
	}
	a.reset()
	return ap
}

func (a *appBuilder) reset() {
	a.ApplicationOptions = options.ApplicationOptions{}
	a.components = nil
	a.shutdown = nil
}

func NewBuilder(opts ...options.ApplicationOptions) AppBuilder {
	opt := options.DefaultApplicationOptions()
	for _, o := range opts {
		opt.Merge(&o)
	}
	return &appBuilder{
		events:             handlers.NullEvents{},
		ApplicationOptions: opt,
	}
}
