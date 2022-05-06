package invoker

import (
	"context"
	"log"
	"reflect"
	"sync"
)

var typeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()
var typeOfWaitGroup = reflect.TypeOf((*sync.WaitGroup)(nil))

type reflectiveInvoke struct {
	fn         reflect.Value
	wgIndex    int
	ctxIndex   int
	inputIndex int
	argsLen    int
	inputType  reflect.Type
}

func (h *reflectiveInvoke) call(ctx context.Context, wg *sync.WaitGroup, input reflect.Value) error {
	args := make([]reflect.Value, h.argsLen)
	if h.ctxIndex >= 0 {
		args[h.ctxIndex] = reflect.ValueOf(ctx)
	}
	if h.wgIndex >= 0 {
		args[h.wgIndex] = reflect.ValueOf(wg)
		wg.Add(1)
	}
	if h.inputIndex >= 0 {
		if input.IsNil() {
			args[h.inputIndex] = input
		} else {
			argInputType := input.Type()
			if !argInputType.AssignableTo(h.inputType) {
				log.Panicf("input type %s is not assignable to %s", argInputType, h.inputType)
			}
			args[h.inputIndex] = input
		}
	}
	result := h.fn.Call(args)
	if len(result) == 0 {
		return nil
	}
	last := result[len(result)-1]
	if !last.IsNil() {
		errRet, ok := last.Interface().(error)
		if ok {
			return errRet
		}
		return nil
	}
	return nil
}

func newReflectiveInvoke(rawFn interface{}) *reflectiveInvoke {
	fn := reflect.ValueOf(rawFn)
	if fn.Kind() != reflect.Func {
		panic("handler must be a function")
	}
	ctxIndex := -1
	inputIndex := -1
	wgIndex := -1
	for i := 0; i < fn.Type().NumIn(); i++ {
		if fn.Type().In(i).AssignableTo(typeOfContext) {
			ctxIndex = i
		} else if fn.Type().In(i).AssignableTo(typeOfWaitGroup) {
			wgIndex = i
		} else {
			inputIndex = i
		}
	}
	argsLen := 0
	if ctxIndex >= 0 {
		argsLen++
	}
	if wgIndex >= 0 {
		argsLen++
	}
	if inputIndex >= 0 {
		argsLen++
	}
	var inputType reflect.Type
	if inputIndex >= 0 {
		inputType = fn.Type().In(inputIndex)
	}
	return &reflectiveInvoke{
		fn:         fn,
		wgIndex:    wgIndex,
		ctxIndex:   ctxIndex,
		inputIndex: inputIndex,
		argsLen:    argsLen,
		inputType:  inputType,
	}
}
