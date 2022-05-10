package apperr

import (
	"runtime"
)

type AppError struct {
	Message    string
	Class      string
	Attributes map[string]interface{}
	Stack      []uintptr
	Cause      error
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func (e *AppError) Error() string {
	switch {
	case e.Message != "" && e.Cause != nil:
		return e.Message + ": " + e.Cause.Error()
	case e.Message != "":
		return e.Message
	case e.Cause != nil:
		return e.Cause.Error()
	default:
		return "unknown error"
	}
}

func (e *AppError) ErrorClass() string { return e.Class }

func (e *AppError) ErrorAttributes() map[string]interface{} { return e.Attributes }

func (e *AppError) StackTrace() []uintptr { return e.Stack }

func newStackTrace(skip int) []uintptr {
	callers := make([]uintptr, 100)
	written := runtime.Callers(skip, callers)
	return callers[:written]
}
