package app

import (
	"context"
	"github.com/yuridevx/app/handlers"
	"github.com/yuridevx/app/options"
	"go.uber.org/atomic"
	"sync"
)

type app struct {
	opts           options.ApplicationOptions
	events         handlers.Events
	components     []*component
	shutdown       []ShutdownFn
	shutdownSignal chan struct{}
	closed         atomic.Bool
}

func (a *app) Run(ctx context.Context, wg *sync.WaitGroup) {
	for _, c := range a.components {
		wg.Add(1)
		go c.goRun(ctx, wg)
	}
	select {
	case <-a.shutdownSignal:
	case <-ctx.Done():
	}
	componentShutdownCtx, cancel := context.WithTimeout(context.Background(), a.opts.ComponentShutdownTimeout)
	for _, c := range a.components {
		c.signalShutdown()
	}
	for _, c := range a.components {
		c.waitExit(componentShutdownCtx)
	}
	cancel()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.opts.AppShutdownTimeout)
	for i := len(a.shutdown) - 1; i >= 0; i-- {
		a.shutdown[i](shutdownCtx)
	}
	cancel()
}

func (a *app) SignalShutdown() {
	if a.closed.Load() {
		return
	}
	close(a.shutdownSignal)
}

var _ Application = (*app)(nil)
