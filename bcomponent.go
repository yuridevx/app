package app

import (
	"github.com/yuridevx/app/handlers"
	"github.com/yuridevx/app/options"
	"reflect"
	"time"
)

type componentBuilder struct {
	options.ComponentOptions
	start     []options.ComponentStartOptions
	shutdown  []options.ComponentShutdownOptions
	cPeriods  []options.CPeriodOptions
	cConsume  []options.CConsumeOptions
	pConsume  []options.PConsumeOptions
	pBlocking []options.PBlockingOptions
}

func (c *componentBuilder) CPeriodIndexFn(
	pIFn options.PeriodIndexFn,
	fn options.HandlerFn,
	opts ...options.CPeriodOptions,
) ComponentBuilder {
	opt := options.DefaultCPeriodOptions()
	opt.PeriodIndexFn = pIFn
	opt.Handler = fn
	for _, o := range opts {
		opt.Merge(&o)
	}
	opt.Validate()
	c.cPeriods = append(c.cPeriods, opt)
	return c
}

func (c *componentBuilder) Options(opts options.ComponentOptions) ComponentBuilder {
	c.ComponentOptions.Merge(&opts)
	return c
}

func (c *componentBuilder) PConsume(
	ch interface{},
	goroutines int,
	fn options.HandlerFn,
	opts ...options.PConsumeOptions,
) ComponentBuilder {
	opt := options.DefaultPConsumeOptions()
	opt.ConsumeCH = ch
	opt.Goroutines = goroutines
	opt.Handler = fn
	for _, o := range opts {
		opt.Merge(&o)
	}
	opt.Validate()
	c.pConsume = append(c.pConsume, opt)
	return c
}

func (c *componentBuilder) PBlocking(fn options.HandlerFn, opts ...options.PBlockingOptions) ComponentBuilder {
	opt := options.DefaultPBlockingOptions()
	opt.Handler = fn
	for _, o := range opts {
		opt.Merge(&o)
	}
	opt.Validate()
	c.pBlocking = append(c.pBlocking, opt)
	return c
}

func (c *componentBuilder) CConsume(
	ch interface{},
	fn options.HandlerFn,
	opts ...options.CConsumeOptions,
) ComponentBuilder {
	opt := options.DefaultCConsumeOptions()
	opt.ConsumeCH = ch
	opt.Handler = fn
	for _, o := range opts {
		opt.Merge(&o)
	}
	opt.Validate()
	c.cConsume = append(c.cConsume, opt)
	return c
}

func (c *componentBuilder) CPeriod(
	p time.Duration,
	fn options.HandlerFn,
	opts ...options.CPeriodOptions,
) ComponentBuilder {
	opt := options.DefaultCPeriodOptions()
	opt.Period = p
	opt.Handler = fn
	for _, o := range opts {
		opt.Merge(&o)
	}
	opt.Validate()
	c.cPeriods = append(c.cPeriods, opt)
	return c
}

func (c *componentBuilder) CPeriodFn(
	pFn options.PeriodFn,
	fn options.HandlerFn,
	opts ...options.CPeriodOptions,
) ComponentBuilder {
	opt := options.DefaultCPeriodOptions()
	opt.PeriodFn = pFn
	opt.Handler = fn
	for _, o := range opts {
		opt.Merge(&o)
	}
	opt.Validate()
	c.cPeriods = append(c.cPeriods, opt)
	return c
}

func (c *componentBuilder) OnStart(fn options.HandlerFn, opts ...options.ComponentStartOptions) ComponentBuilder {
	opt := options.DefaultComponentStartOptions()
	opt.Handler = fn
	for _, o := range opts {
		opt.Merge(&o)
	}
	opt.Validate()
	c.start = append(c.start, opt)
	return c
}

func (c *componentBuilder) OnShutdown(fn options.HandlerFn, opts ...options.ComponentShutdownOptions) ComponentBuilder {
	opt := options.DefaultComponentShutdownOptions()
	opt.Handler = fn
	for _, o := range opts {
		opt.Merge(&o)
	}
	opt.Validate()
	c.shutdown = append(c.shutdown, opt)
	return c
}

func (c *componentBuilder) build(opts options.ApplicationOptions, events handlers.Events) *component {
	start := make([]*handlers.CStartHandler, len(c.start))
	for i, s := range c.start {
		start[i] = handlers.NewStartHandler(handlers.CStart{
			Handler: handlers.Handler{
				HandlerID: reflect.ValueOf(s.Handler).Pointer(),
				Events:    events,
				App:       opts,
				Component: c.ComponentOptions,
			},
			Start: s,
		})
	}
	shutdown := make([]*handlers.CShutdownHandler, len(c.shutdown))
	for i, s := range c.shutdown {
		shutdown[i] = handlers.NewShutdownHandler(handlers.CShutdown{
			Handler: handlers.Handler{
				HandlerID: reflect.ValueOf(s.Handler).Pointer(),
				Events:    events,
				App:       opts,
				Component: c.ComponentOptions,
			},
			Shutdown: s,
		})
	}
	cPeriods := make([]*handlers.CPeriodHandler, len(c.cPeriods))
	for i, p := range c.cPeriods {
		cPeriods[i] = handlers.NewCPeriodHandler(handlers.CPeriod{
			Handler: handlers.Handler{
				HandlerID: reflect.ValueOf(p.Handler).Pointer(),
				Events:    events,
				App:       opts,
				Component: c.ComponentOptions,
			},
			Period: p,
		})
	}
	cConsume := make([]*handlers.CConsumeHandler, len(c.cConsume))
	for i, cc := range c.cConsume {
		cConsume[i] = handlers.NewConsumeHandler(handlers.CConsume{
			Handler: handlers.Handler{
				HandlerID: reflect.ValueOf(cc.Handler).Pointer(),
				Events:    events,
				App:       opts,
				Component: c.ComponentOptions,
			},
			Consume: cc,
		})
	}
	pConsume := make([]*handlers.PConsumeHandler, len(c.pConsume))
	for i, pc := range c.pConsume {
		pConsume[i] = handlers.NewPConsumeHandler(handlers.PConsume{
			Handler: handlers.Handler{
				HandlerID: reflect.ValueOf(pc.Handler).Pointer(),
				Events:    events,
				App:       opts,
				Component: c.ComponentOptions,
			},
			Consume: pc,
		})
	}
	pBlocking := make([]*handlers.PBlockingHandler, len(c.pBlocking))
	for i, pb := range c.pBlocking {
		pBlocking[i] = handlers.NewPBlockingHandler(handlers.PBlocking{
			Handler: handlers.Handler{
				HandlerID: reflect.ValueOf(pb.Handler).Pointer(),
				Events:    events,
				App:       opts,
				Component: c.ComponentOptions,
			},
			Blocking: pb,
		})
	}
	com := &component{
		options:    c.ComponentOptions,
		appOptions: opts,
		start:      start,
		shutdown:   shutdown,
		cPeriod:    cPeriods,
		cConsume:   cConsume,
		pConsume:   pConsume,
		pBlocking:  pBlocking,
		shutdownCh: make(chan struct{}),
		exitCh:     make(chan struct{}),
	}
	c.reset()
	return com
}

func (c *componentBuilder) reset() {
	c.ComponentOptions = options.DefaultComponentOptions()
	c.start = nil
	c.shutdown = nil
	c.cPeriods = nil
	c.cConsume = nil
	c.pConsume = nil
	c.pBlocking = nil
}

var _ ComponentBuilder = (*componentBuilder)(nil)

func newComponentBuilder(def options.ComponentDefinition, opts ...options.ComponentOptions) *componentBuilder {
	opt := options.DefaultComponentOptions()
	opt.Definition = def
	for _, o := range opts {
		opt.Merge(&o)
	}
	return &componentBuilder{
		ComponentOptions: opt,
	}
}
