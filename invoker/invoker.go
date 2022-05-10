package invoker

import (
	"context"
	"github.com/yuridevx/app/options"
	"reflect"
	"sync"
)

type quickContext = func(ctx context.Context, input interface{})
type quickContextErr = func(ctx context.Context, input interface{}) error
type quickContextWg = func(ctx context.Context, wg *sync.WaitGroup, input interface{}) error
type quickSimpleErr = func(ctx context.Context) error
type quickSimple = func(ctx context.Context)
type quickPlain = func()
type quickPlainErr = func() error
type quickInput = func(input interface{})
type quickInputErr = func(input interface{}) error

type Invoker struct {
	qContext    quickContext
	qContextErr quickContextErr
	qContextWg  quickContextWg
	qSimpleErr  quickSimpleErr
	qSimple     quickSimple
	qPlain      quickPlain
	qPlainErr   quickPlainErr
	qInput      quickInput
	qInputErr   quickInputErr

	reflective *reflectiveInvoke
	middleware []options.Middleware
}

func (i *Invoker) Invoke(
	ctx context.Context,
	wg *sync.WaitGroup,
	input interface{},
	call options.Call,
) error {
	nextFn := func(ctx context.Context, input interface{}) error {
		if i.qContext != nil {
			i.qContext(ctx, input)
			return nil
		}

		if i.qContextErr != nil {
			return i.qContextErr(ctx, input)
		}

		if i.qContextWg != nil {
			wg.Add(1)
			return i.qContextWg(ctx, wg, input)
		}

		if i.qSimpleErr != nil {
			return i.qSimpleErr(ctx)
		}

		if i.qSimple != nil {
			i.qSimple(ctx)
			return nil
		}

		if i.qPlain != nil {
			i.qPlain()
			return nil
		}

		if i.qPlainErr != nil {
			return i.qPlainErr()
		}

		if i.qInput != nil {
			i.qInput(input)
			return nil
		}

		if i.qInputErr != nil {
			return i.qInputErr(input)
		}

		return i.reflective.call(ctx, wg, reflect.ValueOf(input))
	}
	err := runMiddleware(ctx, input, call, i.middleware, nextFn)
	return err
}

func (i *Invoker) IsCast() bool {
	return i.qContext != nil || i.qContextErr != nil || i.qContextWg != nil ||
		i.qSimpleErr != nil || i.qSimple != nil || i.qPlain != nil ||
		i.qPlainErr != nil || i.qInput != nil || i.qInputErr != nil
}

func NewInvoker(
	handlerFn options.HandlerFn,
	middleware ...[]options.Middleware,
) *Invoker {
	flat := make([]options.Middleware, 0)
	for _, m := range middleware {
		flat = append(flat, m...)
	}

	i := &Invoker{
		middleware: flat,
	}
	i.qContext, _ = handlerFn.(quickContext)
	i.qContextErr, _ = handlerFn.(quickContextErr)
	i.qContextWg, _ = handlerFn.(quickContextWg)
	i.qSimpleErr, _ = handlerFn.(quickSimpleErr)
	i.qSimple, _ = handlerFn.(quickSimple)
	i.qPlain, _ = handlerFn.(quickPlain)
	i.qPlainErr, _ = handlerFn.(quickPlainErr)
	i.qInput, _ = handlerFn.(quickInput)
	i.qInputErr, _ = handlerFn.(quickInputErr)

	if !i.IsCast() {
		i.reflective = newReflectiveInvoke(handlerFn)
	}

	return i
}
