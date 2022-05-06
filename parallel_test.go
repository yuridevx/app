package app

import (
	"context"
	"go.uber.org/atomic"
	"sync"
	"testing"
	"time"
)

func TestDaemon(t *testing.T) {
	wg := &sync.WaitGroup{}
	ap := NewBuilder()
	res := atomic.NewInt64(0)
	ctx, cancel := context.WithCancel(context.Background())
	ap.C("test").
		CPeriod(time.Millisecond*100, func() {
			cancel()
		}).
		PBlocking(func(ctx context.Context, val interface{}) {
			if val != nil {
				panic("val is not nil")
			}
			for {
				select {
				case <-ctx.Done():
					return
				default:
					res.Inc()
					time.Sleep(20 * time.Millisecond)
				}
			}
		})
	ap.Build().Run(ctx, wg)
	wg.Wait()

	if res.Load() != 5 {
		t.Errorf("res is %d", res.Load())
	}
}

func TestParallel(t *testing.T) {
	wg := &sync.WaitGroup{}
	ap := NewBuilder()
	inp1 := make(chan int, 10)
	inp1 <- 1
	inp1 <- 2
	inp1 <- 3
	inp1 <- 4
	inp1 <- 5
	res := atomic.NewInt64(0)
	ctx, cancel := context.WithCancel(context.Background())
	ap.C("test").
		CPeriod(time.Millisecond*140, func() {
			cancel()
		}).
		PConsume(inp1, 2, func(ctx context.Context, val interface{}) {
			res.Inc()
			time.Sleep(100 * time.Millisecond)
		})
	ap.Build().Run(ctx, wg)
	wg.Wait()

	if res.Load() != 4 {
		t.Error("Expected 4, got ", res.Load())
	}
}
