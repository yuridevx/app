package app

import (
	"context"
	"github.com/yuridevx/app/extension"
	"github.com/yuridevx/app/handlers"
	"github.com/yuridevx/app/invoker"
	"github.com/yuridevx/app/options"
	"go.uber.org/atomic"
	"go.uber.org/multierr"
	"sync"
)

type component struct {
	options    options.ComponentOptions
	appOptions options.ApplicationOptions
	start      []*handlers.CStartHandler
	shutdown   []*handlers.CShutdownHandler
	cPeriod    []*handlers.CPeriodHandler
	cConsume   []*handlers.CConsumeHandler
	pConsume   []*handlers.PConsumeHandler
	pBlocking  []*handlers.PBlockingHandler
	shutdownCh chan struct{}
	closed     atomic.Bool

	exitCh chan struct{}
}

func (c *component) signalShutdown() {
	if c.closed.Load() {
		return
	}
	c.closed.Store(true)
	close(c.shutdownCh)
}

func (c *component) waitExit(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-c.exitCh:
		return
	}
}

func (c *component) goRun(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(c.exitCh)
	defer c.signalShutdown()
	err := c.invokeStart(ctx, wg)
	if err != nil {
		return
	}
	c.launchParallel(ctx, wg)
	c.loop(ctx, wg)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), c.appOptions.ComponentShutdownTimeout)
	c.invokeShutdown(shutdownCtx, wg)
	cancel()
}

func (c *component) launchParallel(ctx context.Context, wg *sync.WaitGroup) {
	for _, p := range c.pBlocking {
		wg.Add(1)
		go p.GoRun(ctx, wg)
	}

	for _, p := range c.pConsume {
		wg.Add(1)
		go p.GoRun(ctx, wg)
	}
}

func (c *component) invokeStart(ctx context.Context, wg *sync.WaitGroup) error {
	if len(c.start) == 0 {
		return nil
	}
	err := invoker.NewInvoker(
		func(ctx context.Context, wg *sync.WaitGroup, input interface{}) error {
			defer wg.Done()
			var allErr error
			for _, handler := range c.start {
				err := handler.Run(ctx, wg)
				if err != nil {
					allErr = multierr.Append(allErr, err)
				}
			}
			return allErr
		}, extension.CallStartGroup,
		c.appOptions.GlobalMiddleware,
		c.options.ComponentMiddleware,
	).Invoke(ctx, wg, nil, nil)
	return err
}

func (c *component) invokeShutdown(ctx context.Context, wg *sync.WaitGroup) {
	if len(c.shutdown) == 0 {
		return
	}
	_ = invoker.NewInvoker(
		func(ctx context.Context, wg *sync.WaitGroup, input interface{}) error {
			defer wg.Done()
			var allErr error
			for _, handler := range c.shutdown {
				err := handler.Run(ctx, wg)
				if err != nil {
					allErr = multierr.Append(allErr, err)
				}
			}
			return allErr
		}, extension.CallShutdownGroup,
		c.appOptions.GlobalMiddleware,
		c.options.ComponentMiddleware,
	).Invoke(ctx, wg, nil, nil)
}

func (c *component) loop(ctx context.Context, wg *sync.WaitGroup) {
	var compete []handlers.CompeteHandler
	for _, cc := range c.cConsume {
		wg.Add(1)
		go cc.GoRun(ctx, wg)
		compete = append(compete, cc)
	}
	for _, cp := range c.cPeriod {
		wg.Add(1)
		go cp.GoRun(ctx, wg)
		compete = append(compete, cp)
	}
	channels := make([]chan interface{}, len(compete))
	for i, h := range compete {
		channels[i] = h.GetSendCh()
	}
	input := handlers.Merge(ctx, channels...)
	for {
		select {
		case <-c.shutdownCh:
			return
		case <-ctx.Done():
			return
		case iv := <-input:
			compete[iv.Index].Execute(ctx, wg, iv.Value)
		}
	}
}
