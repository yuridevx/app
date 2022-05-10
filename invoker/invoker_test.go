package invoker

import (
	"context"
	"sync"
	"testing"
)

type aliasVal = testVal

type testVal struct {
	id string
}

func TestNilInvoke(t *testing.T) {
	called := false
	invoker := NewInvoker(
		func(val *testVal) {
			if val != nil {
				t.Fatal("val is not nil")
			}
			called = true
		},
	)
	wg := &sync.WaitGroup{}

	err := invoker.Invoke(context.Background(), wg, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("called is false")
	}
}

func TestAbstractInvoke(t *testing.T) {
	called := false
	invoker := NewInvoker(
		func(val *aliasVal) {
			if val.id != "test" {
				t.Fatal("val.id is not test")
			}
			called = true
		},
	)
	wg := &sync.WaitGroup{}

	err := invoker.Invoke(context.Background(), wg, &testVal{id: "test"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("called is false")
	}
}
