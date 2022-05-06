package handlers

import "github.com/yuridevx/app/options"

type Handler struct {
	HandlerID uintptr
	Events    Events
	App       options.ApplicationOptions
	Component options.ComponentOptions
}

func (h Handler) GetHandler() uintptr {
	return h.HandlerID
}

func (h Handler) GetComponentDefinition() options.ComponentDefinition {
	return h.Component.Definition
}

type CStart struct {
	Handler
	Start options.ComponentStartOptions
}

type CShutdown struct {
	Handler
	Shutdown options.ComponentShutdownOptions
}

type CPeriod struct {
	Handler
	Period options.CPeriodOptions
}

type CConsume struct {
	Handler
	Consume options.CConsumeOptions
}

type PConsume struct {
	Handler
	Consume options.PConsumeOptions
}

type PBlocking struct {
	Handler
	Blocking options.PBlockingOptions
}
