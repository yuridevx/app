package apperr

type stackTracer interface {
	StackTrace() []uintptr
}

type errorAttributer interface {
	ErrorAttributes() map[string]interface{}
}

func Error(class string, topOptions ...Option) func(args ...interface{}) error {
	return func(args ...interface{}) error {
		err := &AppError{
			Stack: newStackTrace(3),
			Class: class,
		}
		for _, opt := range topOptions {
			opt(err)
		}
		for _, arg := range args {
			switch arg := arg.(type) {
			case error:
				err.Cause = arg
			case string:
				err.Message = arg
			case Option:
				arg(err)
			}
		}
		return err
	}
}
