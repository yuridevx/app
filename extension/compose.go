package extension

import (
	"context"
)

func RunMiddleware(
	ctx context.Context,
	call CallType,
	part Part,
	input interface{},
	middleware []Middleware,
	next NextFn,
) error {
	if len(middleware) == 0 {
		return next(ctx, input)
	}

	index := 0
	var mNext NextFn
	mNext = func(nctx context.Context, ninput interface{}) error {
		if index == len(middleware) {
			return next(nctx, ninput)
		}
		index++
		err := middleware[index-1](nctx, call, ninput, part, mNext)
		return err
	}
	return mNext(ctx, input)
}
