package handlers

import (
	"context"
	"github.com/yuridevx/app/invoker"
	"github.com/yuridevx/app/options"
	"sync"
)

type CStartHandler struct {
	CStart
}

func (h *CStartHandler) Run(ctx context.Context, wg *sync.WaitGroup) error {
	invoke := invoker.NewInvoker(
		h.Start.Handler,
		h.App.Middleware,
		h.Component.Middleware,
		h.Start.Middleware,
	)

	h.Events.Start(h.CStart)
	err := invoke.Invoke(ctx, wg, nil, h.ToCall(options.CallStart))
	h.Events.StartResult(h.CStart, err)
	return err
}

func NewStartHandler(start CStart) *CStartHandler {
	return &CStartHandler{
		CStart: start,
	}
}
