package handlers

import (
	"context"
	"github.com/yuridevx/app/extension"
	"github.com/yuridevx/app/invoker"
	"sync"
)

type CStartHandler struct {
	CStart
}

func (h *CStartHandler) Run(ctx context.Context, wg *sync.WaitGroup) error {
	invoke := invoker.NewInvoker(
		h.Start.Handler,
		extension.CallStart,
		h.App.GlobalMiddleware,
		h.Component.ComponentMiddleware,
		h.Start.CallMiddleware,
	)

	h.Events.Start(h.CStart)
	err := invoke.Invoke(ctx, wg, nil, h.CStart)
	h.Events.StartResult(h.CStart, err)
	return err
}

func NewStartHandler(start CStart) *CStartHandler {
	return &CStartHandler{
		CStart: start,
	}
}
