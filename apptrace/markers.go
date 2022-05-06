package apptrace

import (
	"context"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func Ignore(ctx context.Context) {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return
	}
	txn.Ignore()
}

func Notice(ctx context.Context, err error) {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return
	}
	txn.NoticeError(err)
}
