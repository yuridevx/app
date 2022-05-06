package handlers

import (
	"context"
	"github.com/yuridevx/app/extension"
	"github.com/yuridevx/app/invoker"
	"sync"
)

type CShutdownHandler struct {
	CShutdown
}

func (h *CShutdownHandler) Run(ctx context.Context, wg *sync.WaitGroup) error {
	invoke := invoker.NewInvoker(
		h.Shutdown.Handler,
		extension.CallShutdown,
		h.App.GlobalMiddleware,
		h.Component.ComponentMiddleware,
		h.Shutdown.CallMiddleware,
	)

	h.Events.Shutdown(h.CShutdown)
	err := invoke.Invoke(ctx, wg, nil, h.CShutdown)
	h.Events.ShutdownResult(h.CShutdown, err)
	return err
}

func NewShutdownHandler(shutdown CShutdown) *CShutdownHandler {
	return &CShutdownHandler{
		CShutdown: shutdown,
	}
}
