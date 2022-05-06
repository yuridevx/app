package handlers

import (
	"context"
	"github.com/yuridevx/app/extension"
	"github.com/yuridevx/app/invoker"
	"sync"
)

type PBlockingHandler struct {
	PBlocking
}

func (h *PBlockingHandler) GoRun(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	invoke := invoker.NewInvoker(
		h.Blocking.Handler,
		extension.CallPBlocking,
		h.App.GlobalMiddleware,
		h.Component.ComponentMiddleware,
		h.Blocking.CallMiddleware,
	)

	h.Events.PBlockingStart(h.PBlocking)
	err := invoke.Invoke(ctx, wg, nil, h.PBlocking)
	h.Events.PBlockingResult(h.PBlocking, err)
}

func NewPBlockingHandler(pblocking PBlocking) *PBlockingHandler {
	return &PBlockingHandler{pblocking}
}
