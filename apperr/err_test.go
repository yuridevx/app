package apperr

import (
	"errors"
	"runtime"
	"testing"
)

var TestErr = Error("TestError")

func stackEquals(t *testing.T, left, right uintptr) {
	t.Helper()
	leftFn := runtime.FuncForPC(left)
	rightFn := runtime.FuncForPC(right)
	if leftFn == nil || rightFn == nil {
		t.Errorf("FuncForPC(%d) returned nil", left)
	}
	if leftFn.Name() != rightFn.Name() {
		t.Errorf("FuncForPC(%d) returned %s, expected %s", left, leftFn.Name(), rightFn.Name())
	}
}

func TestError(t *testing.T) {
	ptr, _, _, _ := runtime.Caller(0)
	err1 := errors.New("wrapped error")
	err2 := TestErr(err1, "app error")
	var err3 *AppError
	ok := errors.As(err2, &err3)
	if !ok {
		t.Error("expected error to be an AppError")
	}
	if err3.Class != "TestError" {
		t.Error("expected class to be TestError")
	}
	if err3.Message != "app error" {
		t.Error("expected message to be app error")
	}
	if err3.Cause != err1 {
		t.Error("expected cause to be err1")
	}
	if err3.Error() != "app error: wrapped error" {
		t.Error("expected error to be app error: wrapped error")
	}
	stackEquals(t, ptr, err3.Stack[0])
}
