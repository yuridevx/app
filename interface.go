package app

import (
	"context"
	"github.com/yuridevx/app/handlers"
	"github.com/yuridevx/app/options"
	"sync"
	"time"
)

// Misc Functions

type ShutdownFn = func(ctx context.Context)

//ProxyFn can take *sync.WaitGroup, ctx context.Context
// All other values are treated as input where last overrides previous ones
type ProxyFn = func(args ...interface{}) error

type Application interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	SignalShutdown()
}

type Builder interface {
	Options(opts ...options.ApplicationOption) Builder
	Events(events handlers.Events) Builder
	C(def options.ComponentDefinition, opts ...options.ComponentOption) ComponentBuilder
	OnShutdown(fn ...ShutdownFn) Builder
	Build() Application
}

/*ComponentBuilder
Methods starting with P - "Parallel" runs handlers in parallel
Methods starting with C - "Compete" runs handlers in the same goroutine
All C periodic timers and consume channels compete to run
*/
type ComponentBuilder interface {
	Options(opts ...options.ComponentOption) ComponentBuilder
	PConsume(ch chan interface{}, goroutines int, fn options.HandlerFn, opts ...options.PConsumeOption) ComponentBuilder
	PBlocking(fn options.HandlerFn, opts ...options.PBlockingOption) ComponentBuilder
	CConsume(ch chan interface{}, fn options.HandlerFn, opts ...options.CConsumeOption) ComponentBuilder
	CPeriod(p time.Duration, fn options.HandlerFn, opts ...options.CPeriodOption) ComponentBuilder
	CPeriodFn(pFn options.PeriodFn, fn options.HandlerFn, opts ...options.CPeriodOption) ComponentBuilder
	CPeriodIndexFn(pIFn options.PeriodIndexFn, fn options.HandlerFn, opts ...options.CPeriodOption) ComponentBuilder
	OnStart(fn options.HandlerFn, opts ...options.StartOption) ComponentBuilder
	OnShutdown(fn options.HandlerFn, opts ...options.ShutdownOption) ComponentBuilder
	Proxy(fn options.HandlerFn, opts ...options.ProxyOption) ProxyFn
}
