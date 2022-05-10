package options

import (
	"time"
)

type HandlerFn = interface{}
type PeriodFn = func() time.Duration
type PeriodIndexFn = func(i int32) (int32, time.Duration)

//ComponentDefinition This will be stored, so lightweight values are preferred
type ComponentDefinition = interface{}
