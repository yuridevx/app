package handlers

import "github.com/yuridevx/app/options"

type Call struct {
	Component options.ComponentDefinition
	CallType  options.CallType
	Handler   uintptr
}

func (c *Call) GetCallType() options.CallType {
	return c.CallType
}

func (c *Call) IsCallType(callType ...options.CallType) bool {
	for _, v := range callType {
		if c.CallType == v {
			return true
		}
	}
	return false
}

func (c *Call) GetHandler() uintptr {
	return c.Handler
}

func (c *Call) GetComponentDefinition() interface{} {
	return c.Component
}

var _ options.Call = (*Call)(nil)

func NewCall(component options.ComponentDefinition, callType options.CallType, handler uintptr) *Call {
	return &Call{
		Component: component,
		CallType:  callType,
		Handler:   handler,
	}
}
