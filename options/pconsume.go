package options

type PConsumeOptions struct {
	Handler    HandlerFn
	ConsumeCH  chan interface{}
	Goroutines int
	Middleware []Middleware
}

type PConsumeOption func(o *PConsumeOptions)

func DefaultPConsumeOptions() PConsumeOptions {
	return PConsumeOptions{
		Goroutines: 1,
	}
}

func (a *PConsumeOptions) Validate() {
	if a.Handler == nil {
		panic("handler is required")
	}
	if a.ConsumeCH == nil {
		panic("consume ch is required")
	}
	if a.Goroutines <= 0 {
		panic("goroutines must be greater than 0")
	}
}
