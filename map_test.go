package lck_test

import (
	"testing"

	"github.com/reiver/go-lck"
)

func TestMap_NilReceiver_Get(t *testing.T) {
	var m *lck.Map[string, int]

	value, found := m.Get("anything")
	if found {
		t.Error("expected found=false from nil receiver")
	}
	if value != 0 {
		t.Errorf("expected zero value 0 from nil receiver, got %d", value)
	}
}

func TestMap_NilReceiver_Set(t *testing.T) {
	var m *lck.Map[string, int]

	// Must not panic.
	m.Set("k", 1)
}

func TestMap_NilReceiver_Unset(t *testing.T) {
	var m *lck.Map[string, int]

	// Must not panic.
	m.Unset("k")
}

func TestMap_NilReceiver_Len(t *testing.T) {
	var m *lck.Map[string, int]

	if got := m.Len(); got != 0 {
		t.Errorf("expected 0 from nil receiver, got %d", got)
	}
}

func TestMap_NilReceiver_For(t *testing.T) {
	var m *lck.Map[string, int]

	called := false
	m.For(func(k string, v int) {
		called = true
	})

	if called {
		t.Error("For callback must not be invoked on nil receiver")
	}
}
