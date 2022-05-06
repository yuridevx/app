package handlers

import (
	"context"
	"github.com/yuridevx/app/extension"
	"github.com/yuridevx/app/invoker"
	"sync"
)

type PConsumeHandler struct {
	PConsume
}

func (p *PConsumeHandler) GoRun(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	invoke := invoker.NewInvoker(
		p.Consume.Handler,
		extension.CallPConsume,
		p.App.GlobalMiddleware,
		p.Component.ComponentMiddleware,
		p.Consume.CallMiddleware,
	)

	ch := toInterfaceCh(ctx, p.Consume.ConsumeCH)

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
			err := invoke.Invoke(ctx, wg, val, p.PConsume)
			p.Events.PConsumeResult(p.PConsume, err)
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
