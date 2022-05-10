package options

type ProxyOptions struct {
	Handler    HandlerFn
	Middleware []Middleware
}

type ProxyOption func(o *ProxyOptions)

func DefaultProxyOptions() ProxyOptions {
	return ProxyOptions{}
}

func (o *ProxyOptions) Validate() {
	if o.Handler == nil {
		panic("handler is required")
	}
}
