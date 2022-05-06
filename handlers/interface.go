package handlers

import (
	"context"
	"sync"
)

type CompeteHandler interface {
	GetSendCh() chan interface{}
	Execute(ctx context.Context, wg *sync.WaitGroup, input interface{})
}
