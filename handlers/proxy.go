package handlers

import (
	"context"
	"github.com/yuridevx/app/invoker"
	"github.com/yuridevx/app/options"
	"go.uber.org/atomic"
	"sync"
)

type ProxyHandler struct {
	Proxy
	invoke           *invoker.Invoker
	componentStarted atomic.Bool
	ctx              context.Context
	wg               *sync.WaitGroup
}

func (p *ProxyHandler) Building(handler Handler) {
	p.Proxy.Handler = handler
}

func (p *ProxyHandler) ComponentStarted(ctx context.Context, wg *sync.WaitGroup) {
	p.ctx = ctx
	p.wg = wg
	p.componentStarted.Store(true)
}

func (p *ProxyHandler) ProxyFn(args ...interface{}) error {
	if !p.componentStarted.Load() {
		panic("component not started")
	}
	ctx := p.ctx
	var input interface{}
	for _, arg := range args {
		if argCtx, ok := arg.(context.Context); ok {
			ctx = argCtx
		} else if argInput, ok := arg.(interface{}); ok {
			input = argInput
		}
	}
	return p.invoke.Invoke(ctx, p.wg, input, p.ToCall(options.CallProxy))
}

func NewProxyHandler(proxy Proxy) *ProxyHandler {
	return &ProxyHandler{
		Proxy: proxy,
		invoke: invoker.NewInvoker(
			proxy.Proxy.Handler,
			proxy.Handler.App.Middleware,
			proxy.Handler.Component.Middleware,
			proxy.Proxy.Middleware,
		),
	}
}
