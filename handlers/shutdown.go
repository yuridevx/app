package handlers

import (
	"context"
	"github.com/yuridevx/app/invoker"
	"github.com/yuridevx/app/options"
	"sync"
)

type CShutdownHandler struct {
	CShutdown
}

func (h *CShutdownHandler) Run(ctx context.Context, wg *sync.WaitGroup) error {
	invoke := invoker.NewInvoker(
		h.Shutdown.Handler,
		h.App.Middleware,
		h.Component.Middleware,
		h.Shutdown.Middleware,
	)

	h.Events.Shutdown(h.CShutdown)
	err := invoke.Invoke(ctx, wg, nil, h.ToCall(options.CallShutdown))
	h.Events.ShutdownResult(h.CShutdown, err)
	return err
}

func NewShutdownHandler(shutdown CShutdown) *CShutdownHandler {
	return &CShutdownHandler{
		CShutdown: shutdown,
	}
}
