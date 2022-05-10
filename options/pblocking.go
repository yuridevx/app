package options

type PBlockingOptions struct {
	Handler    HandlerFn
	Middleware []Middleware
}

type PBlockingOption func(o *PBlockingOptions)

func DefaultPBlockingOptions() PBlockingOptions {
	return PBlockingOptions{}
}

func (a *PBlockingOptions) Validate() {
	if a.Handler == nil {
		panic("handler is required")
	}
}
