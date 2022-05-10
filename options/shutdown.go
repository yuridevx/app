package options

type ShutdownOptions struct {
	Handler    HandlerFn
	Middleware []Middleware
}

type ShutdownOption func(o *ShutdownOptions)

func DefaultShutdownOptions() ShutdownOptions {
	return ShutdownOptions{}
}

func (a *ShutdownOptions) Validate() {
	if a.Handler == nil {
		panic("handler is required")
	}
}
