package apptrace

import (
	"context"
	"log"
	"sync"
)

type attrKey int

var attrKeyVal attrKey

func AttributesContext(ctx context.Context, m *sync.Map) context.Context {
	return context.WithValue(ctx, attrKeyVal, m)
}

func AttributesFromContext(ctx context.Context) *sync.Map {
	val := ctx.Value(attrKeyVal)
	if val == nil {
		return nil
	}
	return val.(*sync.Map)
}

func AttributesMap(ctx context.Context, m map[string]interface{}) {
	if m == nil {
		return
	}
	attrs := AttributesFromContext(ctx)
	if attrs == nil {
		return
	}
	for k, v := range m {
		attrs.Store(k, v)
	}
}

func Attributes(ctx context.Context, keysAndValues ...interface{}) {
	if len(keysAndValues)%2 != 0 {
		keysAndValues = keysAndValues[:len(keysAndValues)-1]
		log.Println("apptrace: AttributesAdd: odd number of arguments")
	}
	if len(keysAndValues) == 0 {
		return
	}
	m := AttributesFromContext(ctx)
	if m == nil {
		return
	}
	for i := 0; i < len(keysAndValues); i += 2 {
		m.Store(keysAndValues[i], keysAndValues[i+1])
	}
}
