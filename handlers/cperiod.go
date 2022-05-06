package handlers

import (
	"context"
	"github.com/yuridevx/app/extension"
	"github.com/yuridevx/app/invoker"
	"go.uber.org/atomic"
	"log"
	"math"
	"sync"
	"time"
)

type CPeriodHandler struct {
	CPeriod
	timer   *time.Timer
	sendCh  chan interface{}
	nextCh  chan struct{}
	enabled atomic.Bool
	index   atomic.Int32
	invoke  *invoker.Invoker
}

func (h *CPeriodHandler) launchSwitch(ctx context.Context, wg *sync.WaitGroup) bool {
	if h.Period.SwitchCh == nil {
		h.enabled.Store(true)
		return true
	}
LOOP:
	for {
		select {
		case <-ctx.Done():
			return false
		case enabled, ok := <-h.Period.SwitchCh:
			if !ok {
				return false
			}
			if enabled {
				h.enabled.Store(enabled)
				break LOOP
			}
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case enabled, ok := <-h.Period.SwitchCh:
				if !ok {
					h.enabled.Store(false)
					return
				}
				h.enabled.Store(enabled)
			}
		}
	}()
	return true
}

func (h *CPeriodHandler) GoRun(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer h.Events.CPeriodExit(h.CPeriod)
	ok := h.launchSwitch(ctx, wg)
	if !ok {
		return
	}
	h.reschedule()
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-h.timer.C:
			enabled := h.enabled.Load()
			h.Events.CPeriod(h.CPeriod, t, enabled)
			if !enabled {
				h.reschedule()
				continue
			}
			select {
			case h.sendCh <- t:
				select {
				case <-h.nextCh:
					dur := h.reschedule()
					h.Events.CPeriodNext(h.CPeriod, dur)
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}
}

func (h *CPeriodHandler) reschedule() time.Duration {
	if !h.timer.Stop() {
		select {
		case <-h.timer.C:
		default:
		}
	}
	dur := h.getNextDuration()
	h.timer.Reset(dur)
	return dur
}

func (h *CPeriodHandler) getNextDuration() time.Duration {
	if h.Period.PeriodFn != nil {
		return h.Period.PeriodFn()
	}
	if h.Period.PeriodIndexFn != nil {
		i := h.index.Load()
		newI, dur := h.Period.PeriodIndexFn(i)
		h.index.Store(newI)
		return dur
	}
	if h.Period.Period > 0 {
		return h.Period.Period
	}
	log.Printf("CPeriodHandler: Period is not set")
	return time.Duration(math.MaxInt64)
}

func (h *CPeriodHandler) Execute(ctx context.Context, wg *sync.WaitGroup, input interface{}) {
	h.Events.CPeriodExec(h.CPeriod, input)
	err := h.invoke.Invoke(ctx, wg, input, h.CPeriod)
	h.Events.CPeriodResult(h.CPeriod, err)
	select {
	case h.nextCh <- struct{}{}:
	case <-ctx.Done():
	}
}

func (h *CPeriodHandler) GetSendCh() chan interface{} {
	return h.sendCh
}

func NewCPeriodHandler(cperiod CPeriod) *CPeriodHandler {
	invoke := invoker.NewInvoker(
		cperiod.Period.Handler,
		extension.CallPeriodic,
		cperiod.App.GlobalMiddleware,
		cperiod.Component.ComponentMiddleware,
		cperiod.Period.CallMiddleware,
	)
	return &CPeriodHandler{
		CPeriod: cperiod,
		timer:   time.NewTimer(time.Duration(math.MaxInt64)),
		sendCh:  make(chan interface{}),
		nextCh:  make(chan struct{}),
		invoke:  invoke,
	}
}
