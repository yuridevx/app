package apperr

type Option func(*AppError)

func WithClass(class string) Option {
	return func(e *AppError) {
		e.Class = class
	}
}

func WithAttribute(key string, value interface{}) Option {
	return func(e *AppError) {
		if e.Attributes == nil {
			e.Attributes = make(map[string]interface{})
		}
		e.Attributes[key] = value
	}
}

func WithAttributes(attrs map[string]interface{}) Option {
	return func(e *AppError) {
		if e.Attributes == nil {
			e.Attributes = make(map[string]interface{})
		}
		for k, v := range attrs {
			e.Attributes[k] = v
		}
	}
}

func WithMessage(msg string) Option {
	return func(e *AppError) {
		if e.Message == "" {
			e.Message = msg
		} else {
			e.Message = msg + ": " + e.Message
		}
	}
}

func WithCause(cause error) Option {
	return func(e *AppError) {
		e.Cause = cause
	}
}
