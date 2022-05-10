package handlers

import (
	"context"
	"github.com/yuridevx/app/invoker"
	"github.com/yuridevx/app/options"
	"sync"
)

type PBlockingHandler struct {
	PBlocking
}

func (h *PBlockingHandler) GoRun(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	invoke := invoker.NewInvoker(
		h.Blocking.Handler,
		h.App.Middleware,
		h.Component.Middleware,
		h.Blocking.Middleware,
	)

	h.Events.PBlockingStart(h.PBlocking)
	err := invoke.Invoke(ctx, wg, nil, h.ToCall(options.CallPBlocking))
	h.Events.PBlockingResult(h.PBlocking, err)
}

func NewPBlockingHandler(pblocking PBlocking) *PBlockingHandler {
	return &PBlockingHandler{pblocking}
}
