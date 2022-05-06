package apptrace

import (
	"context"
	"github.com/yuridevx/app/extension"
)

func TraceMiddleware(
	ctx context.Context,
	_ extension.CallType,
	input interface{},
	_ extension.Part,
	next extension.NextFn,
) error {
	ctx, _ = TraceContextNew(ctx)
	err := next(ctx, input)
	return err
}
