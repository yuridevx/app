package handlers

import (
	"context"
	"github.com/yuridevx/app/invoker"
	"github.com/yuridevx/app/options"
	"sync"
)

type CConsumeHandler struct {
	CConsume
	sendCh chan interface{}
	nextCh chan struct{}
	invoke *invoker.Invoker
}

func (c *CConsumeHandler) GoRun(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer c.Events.CConsumeExit(c.CConsume)
	consumeCh := c.CConsume.Consume.ConsumeCh
	for {
		select {
		case <-ctx.Done():
			return
		case item, ok := <-consumeCh:
			if !ok {
				return
			}
			c.Events.CConsume(c.CConsume, item)
			select {
			case c.sendCh <- item:
				select {
				case <-c.nextCh:
					c.Events.CConsumeNext(c.CConsume)
					continue
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}
}

func (c *CConsumeHandler) Execute(ctx context.Context, wg *sync.WaitGroup, input interface{}) {
	c.Events.CConsumeExec(c.CConsume, input)
	err := c.invoke.Invoke(ctx, wg, input, c.ToCall(options.CallCConsume))
	c.Events.CConsumeResult(c.CConsume, err)
	select {
	case c.nextCh <- struct{}{}:
	case <-ctx.Done():
	}
}

func (c *CConsumeHandler) GetSendCh() chan interface{} {
	return c.sendCh
}

func NewConsumeHandler(
	consume CConsume,
) *CConsumeHandler {
	invoke := invoker.NewInvoker(
		consume.Consume.Handler,
		consume.App.Middleware,
		consume.Component.Middleware,
		consume.Consume.Middleware,
	)
	return &CConsumeHandler{
		CConsume: consume,
		invoke:   invoke,
		sendCh:   make(chan interface{}),
		nextCh:   make(chan struct{}),
	}
}
