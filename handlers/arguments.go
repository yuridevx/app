package handlers

import "github.com/yuridevx/app/options"

type Handler struct {
	HandlerID uintptr
	Events    Events
	App       options.ApplicationOptions
	Component options.ComponentOptions
}

type CStart struct {
	Handler
	Start options.StartOptions
}

func (h *CStart) ToCall(callType options.CallType) *Call {
	return &Call{
		Component: h.Component,
		CallType:  callType,
		Handler:   h.HandlerID,
	}
}

type CShutdown struct {
	Handler
	Shutdown options.ShutdownOptions
}

func (h *CShutdown) ToCall(callType options.CallType) *Call {
	return &Call{
		Component: h.Component,
		CallType:  callType,
		Handler:   h.HandlerID,
	}
}

type CPeriod struct {
	Handler
	Period options.CPeriodOptions
}

func (h *CPeriod) ToCall(callType options.CallType) *Call {
	return &Call{
		Component: h.Component,
		CallType:  callType,
		Handler:   h.HandlerID,
	}
}

type CConsume struct {
	Handler
	Consume options.CConsumeOptions
}

func (h *CConsume) ToCall(callType options.CallType) *Call {
	return &Call{
		Component: h.Component,
		CallType:  callType,
		Handler:   h.HandlerID,
	}
}

type PConsume struct {
	Handler
	Consume options.PConsumeOptions
}

func (h *PConsume) ToCall(callType options.CallType) *Call {
	return &Call{
		Component: h.Component,
		CallType:  callType,
		Handler:   h.HandlerID,
	}
}

type PBlocking struct {
	Handler
	Blocking options.PBlockingOptions
}

func (h *PBlocking) ToCall(callType options.CallType) *Call {
	return &Call{
		Component: h.Component,
		CallType:  callType,
		Handler:   h.HandlerID,
	}
}

type Proxy struct {
	Handler
	Proxy options.ProxyOptions
}

func (h *Proxy) ToCall(callType options.CallType) *Call {
	return &Call{
		Component: h.Component,
		CallType:  callType,
		Handler:   h.HandlerID,
	}
}
