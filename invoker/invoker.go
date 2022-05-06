package invoker

import (
	"context"
	"github.com/yuridevx/app/extension"
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

type Invoker struct {
	qContext    quickContext
	qContextErr quickContextErr
	qContextWg  quickContextWg
	qSimpleErr  quickSimpleErr
	qSimple     quickSimple
	qPlain      quickPlain
	call        extension.CallType
	reflective  *reflectiveInvoke
	middleware  []extension.Middleware
}

func (i *Invoker) Invoke(
	ctx context.Context,
	wg *sync.WaitGroup,
	input interface{},
	part extension.Part,
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

		return i.reflective.call(ctx, wg, reflect.ValueOf(input))
	}
	_ = extension.RunMiddleware(ctx, i.call, part, input, i.middleware, nextFn)
	return nil
}

func NewInvoker(
	handlerFn options.HandlerFn,
	call extension.CallType,
	middleware ...[]extension.Middleware,
) *Invoker {
	qContext, _ := handlerFn.(quickContext)
	qContextErr, _ := handlerFn.(quickContextErr)
	qContextWg, _ := handlerFn.(quickContextWg)
	qSimpleErr, _ := handlerFn.(quickSimpleErr)
	qSimple, _ := handlerFn.(quickSimple)
	qPlain, _ := handlerFn.(quickPlain)

	var reflective *reflectiveInvoke
	if qContext == nil && qContextWg == nil && qSimpleErr == nil && qSimple == nil && qPlain == nil {
		reflective = newReflectiveInvoke(handlerFn)
	}

	flat := make([]extension.Middleware, 0)
	for _, m := range middleware {
		flat = append(flat, m...)
	}

	return &Invoker{
		qContext:    qContext,
		qContextErr: qContextErr,
		qContextWg:  qContextWg,
		qSimpleErr:  qSimpleErr,
		qSimple:     qSimple,
		qPlain:      qPlain,
		call:        call,
		reflective:  reflective,
		middleware:  flat,
	}
}
