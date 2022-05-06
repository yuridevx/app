package handlers

import (
	"context"
	"github.com/yuridevx/app/extension"
	"github.com/yuridevx/app/invoker"
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
	consumeCh := toInterfaceCh(ctx, c.CConsume.Consume.ConsumeCH)
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
	err := c.invoke.Invoke(ctx, wg, input, c.CConsume)
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
		extension.CallCConsume,
		consume.App.GlobalMiddleware,
		consume.Component.ComponentMiddleware,
		consume.Consume.CallMiddleware,
	)
	return &CConsumeHandler{
		CConsume: consume,
		invoke:   invoke,
		sendCh:   make(chan interface{}),
		nextCh:   make(chan struct{}),
	}
}
