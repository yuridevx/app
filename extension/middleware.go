package extension

import (
	"context"
)

type CallType string

const (
	CallStartGroup    CallType = "CallStartGroup"
	CallStart         CallType = "CallStart"
	CallShutdownGroup CallType = "CallShutdownGroup"
	CallShutdown      CallType = "CallShutdown"
	CallPeriodic      CallType = "CallPeriodic"
	CallCConsume      CallType = "CallCConsume"
	CallPConsume      CallType = "CallPConsume"
	CallPBlocking     CallType = "CallPBlocking"
)

type Part interface {
	GetHandler() uintptr
}

type Middleware func(ctx context.Context, call CallType, input interface{}, part Part, next NextFn) error
type NextFn func(ctx context.Context, input interface{}) error
