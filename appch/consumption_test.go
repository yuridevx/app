package appch

import (
	"reflect"
	"testing"
)

type obj struct {
	i int
}

func BenchmarkChannelConsumption(b *testing.B) {
	ch := make(chan interface{})
	exit := make(chan struct{})
	defer close(exit)
	go func() {
		i := 0
		for {
			select {
			case ch <- obj{
				i: i,
			}:
			case <-exit:
				return
			}
		}
	}()

	b.Run("normal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			select {
			case <-ch:
			case <-exit:
				return
			}
		}
	})

	cases := []reflect.SelectCase{
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		},
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(exit),
		},
	}

	b.Run("reflect", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reflect.Select(cases)
		}
	})
}

func BenchmarkChannelMerge(b *testing.B) {
	exit := make(chan struct{})
	defer close(exit)
	n := 200
	chs := make([]chan interface{}, n)
	for i := 0; i < n; i++ {
		chs[i] = make(chan interface{})
		go func(i int) {
			for {
				select {
				case chs[i] <- obj{
					i: i,
				}:
				case <-exit:
					return
				}
			}
		}(i)
	}

	normalMerge := Merge(0, exit, chs...)

	b.Run("normal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			select {
			case <-normalMerge:
			case <-exit:
				return
			}
		}
	})

	// Was slow
	//reflectMerge := handlers.MergeReflect(context.Background(), chs...)
	//
	//b.Run("reflect", func(b *testing.B) {
	//	for i := 0; i < b.N; i++ {
	//		select {
	//		case <-reflectMerge:
	//		case <-exit:
	//			return
	//		}
	//	}
	//})
}
