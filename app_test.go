package app

import (
	"context"
	"fmt"
	"github.com/yuridevx/app/appch"
	"github.com/yuridevx/app/handlers"
	"github.com/yuridevx/app/options"
	"log"
	"reflect"
	"sync"
	"testing"
	"time"
)

var quickStart = func(i int32) (int32, time.Duration) {
	if i == 0 {
		return i + 1, time.Nanosecond
	}
	return i, time.Millisecond * 100
}

func TestConsume(t *testing.T) {
	wg := &sync.WaitGroup{}
	inp1 := make(chan int, 10)
	inp1 <- 1
	inp1 <- 2
	inp1 <- 3
	cnt := 0
	ap := NewBuilder().Events(handlers.LogEvents{})
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	ap.C("test").
		CConsume(appch.ToInterfaceChan(inp1), func(ctx context.Context, val interface{}) {
			cnt++
		})
	b := ap.Build()
	b.Run(ctx, wg)
	wg.Wait()
	if cnt != 3 {
		t.Errorf("Expected 3, got %d", cnt)
	}
}

func TestPeriod(t *testing.T) {
	wg := &sync.WaitGroup{}
	cnt := 0
	ap := NewBuilder().Events(handlers.LogEvents{})
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	ap.C("test").
		CPeriodIndexFn(func(i int32) (int32, time.Duration) {
			if i == 0 {
				return i + 1, time.Nanosecond
			}
			return i, time.Millisecond * 20
		}, func(ctx context.Context, val interface{}) {
			cnt++
		})
	b := ap.Build()
	b.Run(ctx, wg)
	wg.Wait()
	if cnt != 5 {
		t.Errorf("Expected 5, got %d", cnt)
	}
}

func TestCompete(t *testing.T) {
	wg := &sync.WaitGroup{}
	cnt1 := 0
	cnt2 := 0
	var res1 []int
	var res2 []int
	inp1 := make(chan int, 10)
	inp2 := make(chan int, 10)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	ap := NewBuilder().Events(handlers.LogEvents{})
	ap.C("test").
		CConsume(appch.ToInterfaceChan(inp1), func(ctx context.Context, val int) {
			res1 = append(res1, val)
		}).
		CConsume(appch.ToInterfaceChan(inp2), func(ctx context.Context, val int) {
			res2 = append(res2, val)
		}).
		CPeriodFn(func() time.Duration {
			return time.Millisecond * 100
		}, func(ctx context.Context) {
			cnt1++
			inp1 <- cnt1
		}).
		CPeriodIndexFn(quickStart, func(ctx context.Context) {
			cnt2++
			inp2 <- cnt2
		}).
		CPeriod(time.Millisecond*300, func() {
			cancel()
		})
	b := ap.Build()
	b.Run(ctx, wg)
	wg.Wait()
	if cnt1 != 2 {
		t.Errorf("cnt1 should be 2, but got %d", cnt1)
	}
	if cnt2 != 3 {
		t.Errorf("cnt2 should be 2, but got %d", cnt2)
	}
	if !reflect.DeepEqual(res1, []int{1, 2}) {
		t.Errorf("res1 should be [1, 2], but got %v", res1)
	}
	if !reflect.DeepEqual(res2, []int{1, 2, 3}) {
		t.Errorf("res2 should be [1, 2, 3], but got %v", res2)
	}
}

func TestComplete(t *testing.T) {
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	go func() {
		<-ctx.Done()
		log.Println("done")
	}()
	ap := NewBuilder()
	ap.C("test").
		CPeriod(time.Millisecond*300, func() {
			cancel()
		})
	ap.Build().Run(ctx, wg)
	wg.Wait()
}

func TestEnable(t *testing.T) {
	wg := &sync.WaitGroup{}
	ap := NewBuilder()
	l := ap.C("test")
	enable := make(chan bool)
	go func() {
		time.Sleep(time.Millisecond * 100)
		enable <- true
	}()
	cnt1 := 0
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	l.CPeriodIndexFn(quickStart, func(ctx context.Context) {
		cnt1++
	}, func(o *options.CPeriodOptions) {
		o.SwitchCh = enable
	})
	l.CPeriod(time.Millisecond*350, func(ctx context.Context) {
		cancel()
	})
	ap.Build().Run(ctx, wg)
	wg.Wait()
	if cnt1 != 3 {
		t.Errorf("cnt1 should be 3, but got %d", cnt1)
	}
}

func TestSwitchOnOff(t *testing.T) {
	wg := &sync.WaitGroup{}
	ap := NewBuilder()
	l := ap.C("test")
	enable := make(chan bool)
	go func() {
		time.Sleep(time.Millisecond * 80)
		enable <- true // 1 startup
		time.Sleep(time.Millisecond * 80)
		enable <- false // make 2nd tick skip
		time.Sleep(time.Millisecond * 80)
		enable <- true // 3 rd tick ok +1 = 2
		time.Sleep(time.Millisecond * 80)
		enable <- false // 4 th tick skip
		time.Sleep(time.Millisecond * 80)
		enable <- true
		// should time out before it takes effect
	}()
	cnt1 := 0
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	l.CPeriodIndexFn(quickStart, func(ctx context.Context) {
		cnt1++
	}, func(o *options.CPeriodOptions) {
		o.SwitchCh = enable
	})
	l.CPeriod(time.Millisecond*420, func(ctx context.Context) {
		cancel()
	})
	ap.Build().Run(ctx, wg)
	wg.Wait()
	if cnt1 != 2 {
		t.Errorf("cnt1 should be 2, but got %d", cnt1)
	}
}

func TestProxyPanic(t *testing.T) {
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	var callArg int

	ap := NewBuilder()
	c := ap.C("test").CPeriod(100*time.Millisecond, func() {
		cancel()
	})
	proxyFn := c.Proxy(func(ctx context.Context, wg *sync.WaitGroup, val int) error {
		defer wg.Done()
		callArg = val
		return fmt.Errorf("test")
	})
	shouldPanic(t, proxyFn)
	ap.Build().Run(ctx, wg)
	errVal := proxyFn(10)
	wg.Wait()
	if callArg != 10 {
		t.Errorf("callArg should be 10, but got %d", callArg)
	}
	if errVal.Error() != "test" {
		t.Errorf("errVal should be test, but got %s", errVal.Error())
	}
}

func shouldPanic(t *testing.T, fn func(...interface{}) error) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("should panic")
		}
	}()
	_ = fn()
}
