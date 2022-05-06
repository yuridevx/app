package handlers

import (
	"context"
	"reflect"
)

func toInterfaceCh(
	ctx context.Context,
	ch interface{},
) chan interface{} {
	if ch == nil {
		return nil
	}
	if output, ok := ch.(chan interface{}); ok {
		return output
	}
	output := make(chan interface{})
	cases := []reflect.SelectCase{
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		},
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ctx.Done()),
		},
	}
	go func() {
		for {
			chose, chv, ok := reflect.Select(cases)
			if chose == 0 {
				if !ok {
					close(output)
					return
				}
				select {
				case output <- chv.Interface():
				case <-ctx.Done():
					return
				}
			} else {
				return
			}
		}
	}()
	return output
}

type IValue struct {
	Index int
	Value interface{}
}

func Merge(ctx context.Context, cs ...chan interface{}) chan IValue {
	out := make(chan IValue)
	for i, c := range cs {
		go func(i int, c <-chan interface{}) {
			for {
				select {
				case v, ok := <-c:
					if !ok {
						return
					}
					select {
					case out <- IValue{i, v}:
					case <-ctx.Done():
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}(i, c)
	}
	return out
}

func MergeReflect(ctx context.Context, cs ...chan interface{}) chan IValue {
	cases := make([]reflect.SelectCase, len(cs))
	for i, c := range cs {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(c),
		}
	}
	chanIndex := len(cases) - 1
	cases = append(cases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ctx.Done()),
	})

	out := make(chan IValue)
	go func() {
		for {
			chose, chv, ok := reflect.Select(cases)
			if chose <= chanIndex {
				if !ok {
					close(out)
					cases[chose].Chan = reflect.ValueOf((chan struct{})(nil))
					return
				}
				select {
				case out <- IValue{chose, chv.Interface()}:
				case <-ctx.Done():
					return
				}
			} else {
				return
			}
		}
	}()
	return out
}
