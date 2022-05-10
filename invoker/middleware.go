package invoker

import (
	"context"
	"github.com/yuridevx/app/options"
)

func runMiddleware(
	ctx context.Context,
	input interface{},
	call options.Call,
	middleware []options.Middleware,
	next options.NextFn,
) error {
	if len(middleware) == 0 {
		return next(ctx, input)
	}

	index := 0
	var mNext options.NextFn
	mNext = func(nctx context.Context, ninput interface{}) error {
		if index == len(middleware) {
			return next(nctx, ninput)
		}
		index++
		err := middleware[index-1](nctx, ninput, call, mNext)
		return err
	}
	return mNext(ctx, input)
}
