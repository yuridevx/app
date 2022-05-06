package benchmark

import (
	"context"
	"github.com/yuridevx/app/handlers"
	"reflect"
	"testing"
	"time"
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
		select {
		case <-ch:
		case <-exit:
			return
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
		reflect.Select(cases)
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

	normalMerge := handlers.Merge(context.Background(), chs...)

	b.Run("normal", func(b *testing.B) {
		time.Sleep(time.Nanosecond * 10)
		select {
		case <-normalMerge:
		case <-exit:
			return
		}
	})

	reflectMerge := handlers.MergeReflect(context.Background(), chs...)

	b.Run("reflect", func(b *testing.B) {
		time.Sleep(time.Nanosecond * 10)
		select {
		case <-reflectMerge:
		case <-exit:
			return
		}
	})
}
