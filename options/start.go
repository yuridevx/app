package options

type StartOptions struct {
	Handler    HandlerFn
	Middleware []Middleware
}

type StartOption func(o *StartOptions)

func DefaultStartOptions() StartOptions {
	return StartOptions{}
}

func (a *StartOptions) Validate() {
	if a.Handler == nil {
		panic("Handler is required")
	}
}
