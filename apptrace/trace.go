package apptrace

import (
	"log"
	"sync"
)

type Trace struct {
	Name       string
	Ignore     bool
	Log        bool
	Attributes map[string]interface{}
	Errors     []error
	Mut        sync.Mutex
}

func (t *Trace) WithName(name string) *Trace {
	if t == nil {
		return nil
	}
	t.Mut.Lock()
	t.Name = name
	t.Mut.Unlock()
	return t
}

func (t *Trace) WithLog(log bool) *Trace {
	if t == nil {
		return nil
	}
	t.Mut.Lock()
	t.Log = log
	t.Mut.Unlock()
	return t
}

func (t *Trace) WithIgnore(ignore bool) *Trace {
	if t == nil {
		return nil
	}
	t.Mut.Lock()
	t.Ignore = ignore
	t.Mut.Unlock()
	return t
}

func (t *Trace) WithAttributes(keysAndValues ...interface{}) *Trace {
	if t == nil {
		return nil
	}
	if len(keysAndValues) == 0 {
		return t
	}
	if len(keysAndValues)%2 != 0 {
		log.Println("apptrace: WithAttributes: number of arguments must be even")
		keysAndValues = keysAndValues[:len(keysAndValues)-1]
	}
	t.Mut.Lock()
	defer t.Mut.Unlock()
	for i := 0; i < len(keysAndValues); i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			log.Println("apptrace: invalid key", keysAndValues[i])
			continue
		}
		value := keysAndValues[i+1]
		t.Attributes[key] = value
	}
	return t
}

func (t *Trace) WithAttributesMap(m map[string]interface{}) *Trace {
	if t == nil {
		return nil
	}
	if len(m) == 0 {
		return t
	}
	t.Mut.Lock()
	defer t.Mut.Unlock()
	for k, v := range m {
		t.Attributes[k] = v
	}
	return t
}

func (t *Trace) WithError(err error) *Trace {
	if t == nil {
		return nil
	}
	t.Mut.Lock()
	t.Errors = append(t.Errors, err)
	t.Mut.Unlock()
	return t
}

func (t *Trace) GetAttributes() map[string]interface{} {
	if t == nil {
		return nil
	}
	t.Mut.Lock()
	defer t.Mut.Unlock()
	if len(t.Attributes) == 0 {
		return nil
	}
	copied := make(map[string]interface{})
	for k, v := range t.Attributes {
		copied[k] = v
	}
	return copied
}

func (t *Trace) GetErrors() []error {
	if t == nil {
		return nil
	}
	t.Mut.Lock()
	defer t.Mut.Unlock()
	if len(t.Errors) == 0 {
		return nil
	}
	copied := make([]error, len(t.Errors))
	copy(copied, t.Errors)
	return copied
}

func (t *Trace) GetName() string {
	if t == nil {
		return "unknown trace"
	}
	t.Mut.Lock()
	defer t.Mut.Unlock()
	return t.Name
}

func NewTrace() *Trace {
	return &Trace{
		Attributes: make(map[string]interface{}),
	}
}
