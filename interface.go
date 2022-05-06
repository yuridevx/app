package app

import (
	"context"
	"github.com/yuridevx/app/handlers"
	"github.com/yuridevx/app/options"
	"sync"
	"time"
)

// Misc Functions

type AppShutdownFn func(ctx context.Context)

type Application interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	SignalShutdown()
}

type AppBuilder interface {
	Options(opts options.ApplicationOptions) AppBuilder
	Events(events handlers.Events) AppBuilder
	C(def options.ComponentDefinition, opts ...options.ComponentOptions) ComponentBuilder
	OnShutdown(fn ...AppShutdownFn) AppBuilder
	Build() Application
}

type ComponentBuilder interface {
	Options(opts options.ComponentOptions) ComponentBuilder
	PConsume(ch interface{}, goroutines int, fn options.HandlerFn, opts ...options.PConsumeOptions) ComponentBuilder
	PBlocking(fn options.HandlerFn, opts ...options.PBlockingOptions) ComponentBuilder
	CConsume(ch interface{}, fn options.HandlerFn, opts ...options.CConsumeOptions) ComponentBuilder
	CPeriod(p time.Duration, fn options.HandlerFn, opts ...options.CPeriodOptions) ComponentBuilder
	CPeriodFn(pFn options.PeriodFn, fn options.HandlerFn, opts ...options.CPeriodOptions) ComponentBuilder
	CPeriodIndexFn(pIFn options.PeriodIndexFn, fn options.HandlerFn, opts ...options.CPeriodOptions) ComponentBuilder
	OnStart(fn options.HandlerFn, opts ...options.ComponentStartOptions) ComponentBuilder
	OnShutdown(fn options.HandlerFn, opts ...options.ComponentShutdownOptions) ComponentBuilder
}
