package handlers

import (
	"context"
	"github.com/yuridevx/app/invoker"
	"github.com/yuridevx/app/options"
	"sync"
)

type PConsumeHandler struct {
	PConsume
}

func (p *PConsumeHandler) GoRun(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	invoke := invoker.NewInvoker(
		p.Consume.Handler,
		p.App.Middleware,
		p.Component.Middleware,
		p.Consume.Middleware,
	)

	ch := p.Consume.ConsumeCH

	for i := 0; i < p.Consume.Goroutines; i++ {
		wg.Add(1)
		go p.goLoop(ctx, wg, ch, invoke)
	}
}

func (p *PConsumeHandler) goLoop(
	ctx context.Context,
	wg *sync.WaitGroup,
	ch chan interface{},
	invoke *invoker.Invoker,
) {
	defer wg.Done()
	defer p.Events.PConsumeExit(p.PConsume)
	for {
		select {
		case <-ctx.Done():
			return
		case val, ok := <-ch:
			if !ok {
				return
			}
			p.Events.PConsume(p.PConsume, val)
			err := invoke.Invoke(ctx, wg, val, p.ToCall(options.CallPConsume))
			p.Events.PConsumeResult(p.PConsume, err)
			// more stable timing on context done
			// take a look at TestParallel
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}
}

func NewPConsumeHandler(
	pConsume PConsume,
) *PConsumeHandler {
	return &PConsumeHandler{
		PConsume: pConsume,
	}
}
