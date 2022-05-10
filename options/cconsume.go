package options

type CConsumeOptions struct {
	Handler    HandlerFn
	ConsumeCh  chan interface{}
	Middleware []Middleware
}

type CConsumeOption func(o *CConsumeOptions)

func DefaultCConsumeOptions() CConsumeOptions {
	return CConsumeOptions{}
}

func (a *CConsumeOptions) Validate() {
	if a.Handler == nil {
		panic("Handler is required")
	}
	if a.ConsumeCh == nil {
		panic("ConsumeCh is required")
	}
}
