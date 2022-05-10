package options

import (
	"context"
)

type CallType string

const (
	CallStartGroup    CallType = "CallStartGroup"
	CallStart         CallType = "CallStart"
	CallShutdownGroup CallType = "CallShutdownGroup"
	CallShutdown      CallType = "CallShutdown"
	CallCPeriod       CallType = "CallCPeriod"
	CallCConsume      CallType = "CallCConsume"
	CallPConsume      CallType = "CallPConsume"
	CallPBlocking     CallType = "CallPBlocking"
	CallProxy         CallType = "CallProxy"
)

type Call interface {
	GetCallType() CallType
	IsCallType(callType ...CallType) bool
	GetHandler() uintptr
	GetComponentDefinition() interface{}
}

type Middleware = func(ctx context.Context, input interface{}, call Call, next NextFn) error
type NextFn = func(ctx context.Context, input interface{}) error
