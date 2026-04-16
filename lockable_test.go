package lck_test

import (
	"testing"

	"github.com/reiver/go-lck"
)

func TestLockable_GetZeroValue(t *testing.T) {
	var l lck.Lockable[int]

	if got := l.Get(); got != 0 {
		t.Errorf("expected zero value 0, got %d", got)
	}
}

func TestLockable_SetThenGet(t *testing.T) {
	var l lck.Lockable[int]

	l.Set(42)

	if got := l.Get(); got != 42 {
		t.Errorf("expected 42, got %d", got)
	}
}

func TestLockable_Let_Mutates(t *testing.T) {
	var l lck.Lockable[int]

	l.Set(10)
	l.Let(func(v int) int {
		return v + 5
	})

	if got := l.Get(); got != 15 {
		t.Errorf("expected 15, got %d", got)
	}
}

func TestLockable_Let_OnZeroValue(t *testing.T) {
	var l lck.Lockable[string]

	l.Let(func(v string) string {
		return "hello"
	})

	if got := l.Get(); got != "hello" {
		t.Errorf("expected %q, got %q", "hello", got)
	}
}

func TestLockable_NilReceiver_Get(t *testing.T) {
	var l *lck.Lockable[int]

	if got := l.Get(); got != 0 {
		t.Errorf("expected zero value 0 from nil receiver, got %d", got)
	}
}

func TestLockable_NilReceiver_Set(t *testing.T) {
	var l *lck.Lockable[int]

	// Must not panic.
	l.Set(99)
}

func TestLockable_NilReceiver_Let(t *testing.T) {
	var l *lck.Lockable[int]

	called := false
	l.Let(func(v int) int {
		called = true
		return v
	})

	if called {
		t.Error("Let callback must not be invoked on nil receiver")
	}
}
